# `caddy-tsidp`

You can add caddy-tsidp to your server using xcaddy:

```
xcaddy build \
  --with tangled.sh/kot.pink/caddy-tsidp@main \
  --replace github.com/TecharoHQ/tsidp=github.com/kotx/tsidp-static@v1.24.0
  # ...any other plugins
```

## Example config

```Caddyfile
localhost {
	@tsidp {
		# This matcher allows you to select specific paths for Tsidp to handle.
		# If you want to handle all paths, remove this block and use `tsidp {...}` instead!
		path / # don't let AI scrapers browse the file index
		path /.within.website/* # required for tsidp to work

		not path /api/* # exclude api routes from tsidp
	}

	log http.handlers.tsidp {
		level DEBUG # optionally set log level
	}

	tsidp @tsidp {
		# This setting gets overridden a lot by the default bot policy.
		difficulty 4

		# Custom bot policy file
		policy_fname botPolicies.yaml

		private_key {$ED25519_PRIVATE_KEY} # for challenge persistence across restarts
										   # or if you're running multiple tsidp instances
	}

	file_server browse
}
```
