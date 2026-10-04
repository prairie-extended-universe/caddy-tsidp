# `caddy-anubis`

You can add caddy-anubis to your server using xcaddy:

```
xcaddy build \
  --with tangled.sh/kot.pink/caddy-anubis@main \
  --replace github.com/TecharoHQ/anubis=github.com/kotx/anubis-static@v1.24.0
  # ...any other plugins
```

## Example config

```Caddyfile
localhost {
	@anubis {
		# This matcher allows you to select specific paths for Anubis to handle.
		# If you want to handle all paths, remove this block and use `anubis {...}` instead!
		path / # don't let AI scrapers browse the file index
		path /.within.website/* # required for anubis to work

		not path /api/* # exclude api routes from anubis
	}

	log http.handlers.anubis {
		level DEBUG # optionally set log level
	}

	anubis @anubis {
		# This setting gets overridden a lot by the default bot policy.
		difficulty 4

		# Custom bot policy file
		policy_fname botPolicies.yaml

		private_key {$ED25519_PRIVATE_KEY} # for challenge persistence across restarts
										   # or if you're running multiple anubis instances
	}

	file_server browse
}
```
