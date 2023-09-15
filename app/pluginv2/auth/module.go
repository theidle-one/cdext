package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/caddyserver/caddy/caddyhttp/httpserver"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	caddy.RegisterModule(Middleware{})
	httpcaddyfile.RegisterHandlerDirective("authenticator", parseCaddyfile)
}

type ResponseBody struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    Body   `json:"data,omitempty"`
}

type User struct {
	Username string   `bson:"username" json:"username"`
	Name     string   `bson:"name,omitempty" json:"name,omitempty"`
	Disabled bool     `bson:"disabled" json:"disabled"`
	Scopes   []string `bson:"scopes" json:"scopes"`
}

type License struct {
	Name   string   `bson:"name" json:"name"`
	Scopes []string `bson:"scopes" json:"scopes"`
}

type Body struct {
	LicenseID string  `bson:"licenceID" json:"licenceID"`
	Scope     string  `json:"scope"`
	User      User    `json:"user"`
	License   License `json:"license"`
}

type Rule struct {
	Endpoint  string
	Except    []string
	Resources []string
}

// Middleware implements an HTTP handler that writes the
// visitor's IP address to a file or stream.
type Middleware struct {
	Rules []Rule
}

// CaddyModule returns the Caddy module information.
func (Middleware) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.authenticator",
		New: func() caddy.Module { return new(Middleware) },
	}
}

// Provision implements caddy.Provisioner.
func (m *Middleware) Provision(ctx caddy.Context) error {
	return nil
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	var protected, isAuthenticated bool
	for _, rule := range m.Rules {
		for _, res := range rule.Except {
			parts := strings.Fields(res)
			if len(parts) == 2 {
				if httpserver.Path(r.URL.Path).Matches(parts[0]) && parts[1] == r.Method {
					goto NEXT
				}
			}

			if len(parts) == 1 {
				if httpserver.Path(r.URL.Path).Matches(parts[0]) {
					goto NEXT
				}
			}
		}
	}

	for _, rule := range m.Rules {
		for _, res := range rule.Resources {
			if !httpserver.Path(r.URL.Path).Matches(res) {
				continue
			}

			// path matches; this endpoint is protected
			protected = true

			req, err := http.NewRequest(http.MethodGet, rule.Endpoint, nil)
			if err != nil {
				continue
			}
			req.Header.Set("Authorization", r.Header.Get("Authorization"))

			resp, err := http.DefaultClient.Do(req)
			if err != nil || resp.StatusCode != http.StatusOK {
				continue
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				continue
			}
			resp.Body.Close()

			var respBody ResponseBody
			if err := json.Unmarshal(body, &respBody); err != nil {
				continue
			}

			// by this point, authentication was successful
			isAuthenticated = true
			u := respBody.Data

			fmt.Println(u)

			r.Header.Add("LicenseID", u.LicenseID)
			r = r.WithContext(context.WithValue(r.Context(), "licenceID", u.License))
			r = r.WithContext(context.WithValue(r.Context(), "user.username", u.User.Username))
			r = r.WithContext(context.WithValue(r.Context(), "scope", u.Scope))
			r = r.WithContext(context.WithValue(r.Context(), "user.scopes", u.User.Scopes))
			r = r.WithContext(context.WithValue(r.Context(), "license.scopes", u.License.Scopes))
			r = r.WithContext(context.WithValue(r.Context(), "license.name", u.License.Name))

			goto NEXT
		}
	}

	if protected && !isAuthenticated {
		w.WriteHeader(http.StatusUnauthorized)
		return errors.New("unauthenticated")
	}

NEXT:
	fmt.Printf("[Authenticator Middleware]: %v-%v-%v\n", r.Context().Value("user.licenceID"), r.Context().Value("user.username"), r.Context().Value("user.scopes"))
	return next.ServeHTTP(w, r)
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
func (m *Middleware) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	var rules []Rule
	for d.Next() {
		var rule Rule

		for d.NextBlock(0) {
			val := d.Val()
			args := d.RemainingArgs()

			switch val {
			case "endpoint":
				rule.Endpoint = args[0]
			case "except":
				rule.Except = append(rule.Except, strings.Join(args, " "))
			case "require":
				rule.Resources = append(rule.Resources, args...)
			}
		}

		rules = append(rules, rule)
	}

	m.Rules = rules
	return nil
}

// parseCaddyfile unmarshals tokens from h into a new Middleware.
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	m := new(Middleware)
	err := m.UnmarshalCaddyfile(h.Dispenser)
	return m, err
}

// Interface guards
var (
	_ caddy.Provisioner           = (*Middleware)(nil)
	_ caddyhttp.MiddlewareHandler = (*Middleware)(nil)
	_ caddyfile.Unmarshaler       = (*Middleware)(nil)
)
