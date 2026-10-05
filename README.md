# `caddy-tsidp`

You can add caddy-tsidp to your server using xcaddy:

```
xcaddy build \
  --with github.com/prairie-extended-universe/caddy-tsidp@trunk
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

## Grants and other config

You still need to configure grants per the [tsidp docs](https://github.com/tailscale/tsidp/blob/v0.0.15/README.md). Most of the rest of the config is handled in the Caddyfile.
