# `caddy-tsidp`

You can add caddy-tsidp to your server using xcaddy:

```
xcaddy build \
  --with tangled.sh/kot.pink/caddy-tsidp@main
  # ...any other plugins
```

This assumes caddy is running next to tailscaled, and you're not using `tailscale serve` or other proxying.

## Example config

```Caddyfile
idp.lan.mynet {
	log http.handlers.tsidp {
		level DEBUG # optionally set log level
	}

	tsidp {
		# Required
		state_dir ./idp-state
		# Probably required
		canonical_hostname idp.lan.mynet
		# Completely optional
		enable_sts
	}
}
```
