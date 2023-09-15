package systemLog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/caddyserver/caddy/caddyhttp/httpserver"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	caddy.RegisterModule(Middleware{})
	httpcaddyfile.RegisterHandlerDirective("system_log", parseCaddyfile)
}

type ResponseBody struct {
	Message string  `json:"message"`
	Code    int     `json:"code"`
	Data    Licence `json:"data,omitempty"`
}

type Licence struct {
	ID     string   `bson:"_id" json:"id"`
	Status int      `bson:"status" json:"status"`
	Name   string   `bson:"name" json:"name"`
	Domain string   `bson:"domain" json:"domain"`
	Scopes []string `bson:"scopes" json:"scopes"`
}

type Rule struct {
	Endpoint     string
	Except       []string
	Resources    []string
	ServiceToken string
}

// Middleware implements an HTTP handler that writes the
// visitor's IP address to a file or stream.
type Middleware struct {
	Rules []Rule
}

// CaddyModule returns the Caddy module information.
func (Middleware) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.system_log",
		New: func() caddy.Module { return new(Middleware) },
	}
}

// Provision implements caddy.Provisioner.
func (m *Middleware) Provision(ctx caddy.Context) error {
	return nil
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (m Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	for _, rule := range m.Rules {
		for _, res := range rule.Except {
			if httpserver.Path(r.URL.Path).Matches(res) {
				goto NEXT
			}
		}
	}

	for _, rule := range m.Rules {
		for _, res := range rule.Resources {
			if !httpserver.Path(r.URL.Path).Matches(res) {
				continue
			}

			var (
				data     interface{}
				resource string
				action   string
				service  string
				log      LicenceSystemLog
			)
			resource = "N/A"

			switch r.Method {
			case http.MethodPut, http.MethodPatch:
				action = ACTION_UPDATE
			case http.MethodPost:
				action = ACTION_ADD
			case http.MethodDelete:
				action = ACTION_DELETE
			case http.MethodGet:
				action = ACTION_GET
			}

			paths := strings.Split(r.URL.String(), "/")
			if len(paths) > 4 {
				service = paths[1]
				resource = paths[4]
			}

			if len(paths) > 5 {
				for i, v := range paths {
					if i < 5 {
						continue
					}
					resource = fmt.Sprintf("%s/%s", resource, v)
				}

				path := strings.ToLower(strings.TrimSpace(paths[5]))
				if path == "download" {
					resource = paths[4]
					action = ACTION_DOWNLOAD
				}

				if path == "upload" {
					resource = paths[4]
					action = ACTION_UPLOAD
				}
			}

			bodyBytes, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(bodyBytes, &data)
			r.Body.Close()

			log.Action = action
			log.Resource = strings.Split(resource, "?")[0]
			log.Body = data
			log.CreatedAt = time.Now()
			log.Service = service
			log.Username = r.Context().Value("user.username").(string)
			log.LicenceID = r.Context().Value("licence.id").(string)

			r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

			body, err := json.Marshal(log)
			if err != nil {
				continue
			}

			req, err := http.NewRequest(http.MethodPost, rule.Endpoint, bytes.NewBuffer(body))
			if err != nil {
				continue
			}
			req.Header.Set("Authorization", fmt.Sprintf("Bearer%s", rule.ServiceToken))

			resp, err := http.DefaultClient.Do(req)
			if err != nil || resp.StatusCode != http.StatusOK {
				continue
			}
			goto NEXT
		}
	}

NEXT:
	fmt.Printf("[SystemLog Middleware]: %v-%v-%v\n", r.Context().Value("licence.id"), r.Context().Value("licence.name"), r.Context().Value("licence.scopes"))
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
				rule.Except = append(rule.Except, args...)
			case "require":
				rule.Resources = append(rule.Resources, args...)
			case "serviceToken":
				rule.ServiceToken = args[0]
			}

		}

		rules = append(rules, rule)
	}

	m.Rules = rules
	return nil
}

// parseCaddyfile unmarshals tokens from h into a new Middleware.
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var m Middleware
	err := m.UnmarshalCaddyfile(h.Dispenser)
	return m, err
}

// Interface guards
var (
	_ caddy.Provisioner           = (*Middleware)(nil)
	_ caddyhttp.MiddlewareHandler = (*Middleware)(nil)
	_ caddyfile.Unmarshaler       = (*Middleware)(nil)
)
