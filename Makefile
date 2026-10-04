caddy: *.go go.mod go.sum tsidp/server/*.go
	xcaddy build \
		--with git.gay/astraluma/caddy-tsidp=. \
		--replace github.com/tailscale/tsidp=./tsidp
