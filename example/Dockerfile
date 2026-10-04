FROM caddy:builder AS builder

RUN xcaddy build \
    --with tangled.sh/kot.pink/caddy-anubis \
    --replace github.com/TecharoHQ/anubis=github.com/kotx/anubis-static@v1.24.0

FROM caddy:latest

COPY --from=builder /usr/bin/caddy /usr/bin/caddy
COPY Caddyfile /etc/caddy/Caddyfile
COPY botPolicy.yaml .
