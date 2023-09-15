package authoritative

import (
	"errors"
	"fmt"
	"net/http"

	"git.cyradar.com/microservices/api-gateway/app/pluginv2"

	"github.com/caddyserver/caddy/caddyhttp/httpserver"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	caddy.RegisterModule(Middleware{})
	httpcaddyfile.RegisterHandlerDirective("authoritative", parseCaddyfile)
}

type Rule struct {
	Except    []string
	PathScope map[string]map[string]string
	Paths     []string
}

// Middleware implements an HTTP handler that writes the
// visitor's IP address to a file or stream.
type Middleware struct {
	Rules []Rule
}

// CaddyModule returns the Caddy module information.
func (Middleware) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.authoritative",
		New: func() caddy.Module { return new(Middleware) },
	}
}

// Provision implements caddy.Provisioner.
func (m *Middleware) Provision(ctx caddy.Context) error {
	return nil
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (m Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	var isAuthorized bool

	for _, rule := range m.Rules {
		for _, res := range rule.Except {
			if httpserver.Path(r.URL.Path).Matches(res) {
				goto NEXT
			}
		}
	}

LOOP_RULE:
	for _, rule := range m.Rules {
		for _, res := range rule.Paths {

			if !pluginv2.CompareURI(r, res) {
				continue
			}

			if _, ok := rule.PathScope[res]; !ok {
				continue
			}

			mapMethodScope := rule.PathScope[res]

			for method, scp := range mapMethodScope {
				if method != "*" && r.Method != method {
					continue
				}

				if err := DoScopesAllowed(r, scp); err != nil {
					isAuthorized = false
					break LOOP_RULE
				}

				isAuthorized = true
				continue
			}
		}
	}

	if !isAuthorized {
		w.WriteHeader(http.StatusUnauthorized)
		return errors.New("unauthorized permission")
	}

NEXT:
	fmt.Printf("[Authoritative Middleware]: %v-%v-%v | %v-%v-%v\n",
		r.Context().Value("licenseID"), r.Context().Value("user.username"), r.Context().Value("user.scopes"),
		r.Context().Value("licenseID"), r.Context().Value("license.name"), r.Context().Value("license.scopes"),
	)
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
			case "except":
				rule.Except = append(rule.Except, args...)
			default:
				method := val
				path := args[0]
				scope := args[1]

				if rule.PathScope == nil {
					rule.PathScope = make(map[string]map[string]string)
				}

				if rule.PathScope[path] == nil {
					rule.Paths = append(rule.Paths, path)
					rule.PathScope[path] = make(map[string]string)
				}
				rule.PathScope[path][method] = scope
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

func DoScopesAllowed(r *http.Request, scp string) error {
	if !CheckLicenceScopePermission(r, scp) {
		return errors.New("license authorization error")
	}

	if r.Context().Value("scope") == "user" {
		if !CheckUserScopePermission(r, scp) {
			return errors.New("user authorization error")
		}
	}
	return nil
}

func CheckLicenceScopePermission(r *http.Request, scp string) (result bool) {
	linScp, _ := r.Context().Value("license.scopes").([]string)
	return NewScopes(linScp...).Contains(NewScope(scp))
}

func CheckUserScopePermission(r *http.Request, scp string) (result bool) {
	uScp, _ := r.Context().Value("user.scopes").([]string)
	return NewScopes(uScp...).Contains(NewScope(scp))
}

// Interface guards
var (
	_ caddy.Provisioner           = (*Middleware)(nil)
	_ caddyhttp.MiddlewareHandler = (*Middleware)(nil)
	_ caddyfile.Unmarshaler       = (*Middleware)(nil)
)
