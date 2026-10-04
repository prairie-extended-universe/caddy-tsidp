set windows-powershell := true

# Show this help
@help:
  just --list


build:
  make caddy

run: build
  ./caddy run
