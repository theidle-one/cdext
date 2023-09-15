package main

import (
	caddycmd "github.com/caddyserver/caddy/v2/cmd"

	_ "git.cyradar.com/microservices/api-gateway/app/pluginv2/auth"
	_ "git.cyradar.com/microservices/api-gateway/app/pluginv2/authoritative"
	_ "github.com/caddyserver/caddy/v2/modules/standard"
	_ "github.com/mholt/caddy-ratelimit"
	_ "github.com/sillygod/cdp-cache"
)

func main() {
	caddycmd.Main()
}
