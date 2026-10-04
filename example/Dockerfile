FROM caddy:builder AS builder

RUN xcaddy build \
    --with tangled.sh/kot.pink/caddy-tsidp \
    --replace github.com/TecharoHQ/tsidp=github.com/kotx/tsidp-static@v1.24.0

FROM caddy:latest

COPY --from=builder /usr/bin/caddy /usr/bin/caddy
COPY Caddyfile /etc/caddy/Caddyfile
COPY botPolicy.yaml .
