package pluginv2

import (
	"net/http"
	"strings"

	"github.com/caddyserver/caddy/caddyhttp/httpserver"
)

func CompareURI(r *http.Request, reqpath string) bool {
	if r.Method != "POST" {
		return httpserver.Path(r.URL.Path).Matches(reqpath)
	}
	if len(strings.Split(r.URL.Path, "/")) < len(strings.Split(reqpath, "/")) {
		return false
	}
	if strings.HasSuffix(reqpath, "*") {
		req := r.URL.Path
		if strings.HasSuffix(req, "/") {
			req = strings.TrimSuffix(req, "/")
		}
		buffer := strings.Split(req, "/")
		req = strings.Join(buffer[:len(buffer)-1], "/")
		reqpath = strings.TrimSuffix(reqpath, "/*")
		return httpserver.Path(req).Matches(reqpath)
	}
	return httpserver.Path(r.URL.Path).Matches(reqpath)
}
