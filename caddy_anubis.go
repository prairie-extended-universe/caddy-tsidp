package caddy_anubis

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/TecharoHQ/anubis"
	libanubis "github.com/TecharoHQ/anubis/lib"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(AnubisMiddleware{})
	httpcaddyfile.RegisterHandlerDirective("anubis", parseCaddyfileHandler)
	httpcaddyfile.RegisterDirectiveOrder("anubis", httpcaddyfile.Before, "push")
}

func (AnubisMiddleware) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.anubis",
		New: func() caddy.Module { return new(AnubisMiddleware) },
	}
}

type AnubisMiddleware struct {
	Options           libanubis.Options `json:"options"`
	PolicyFname       string            `json:"policy_fname,omitempty"`
	DefaultDifficulty int               `json:"default_difficulty,omitempty"`

	anubis *libanubis.Server
	log    *zap.Logger
	next   caddyhttp.Handler
	err    error
}

// Interface guards
var (
	_ caddyhttp.MiddlewareHandler = (*AnubisMiddleware)(nil)
	_ caddyfile.Unmarshaler       = (*AnubisMiddleware)(nil)
	_ caddy.Provisioner           = (*AnubisMiddleware)(nil)
)

func (m *AnubisMiddleware) Provision(ctx caddy.Context) error {
	m.log = ctx.Logger()
	m.Options.Logger = ctx.Slogger()

	m.log.Debug("loading anubis policies", zap.String("policy_file", m.PolicyFname), zap.Int("default_difficulty", m.DefaultDifficulty))
	policy, err := libanubis.LoadPoliciesOrDefault(ctx, m.PolicyFname, m.DefaultDifficulty, ctx.Logger().Level().String(), false)
	if err != nil {
		return fmt.Errorf("failed to load anubis policies from '%s': %w", m.PolicyFname, err)
	}

	m.Options.Next = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.err = m.next.ServeHTTP(w, r); err != nil {
			m.log.Debug("received error from next handler", zap.Error(err))
		}
	})
	m.Options.Policy = policy
	m.anubis, err = libanubis.New(m.Options)
	if err != nil {
		return err
	}

	return nil
}

func (m *AnubisMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	remoteHost, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return err
	}
	r.Header.Set("X-Real-Ip", remoteHost)
	r.Header.Set("X-Http-Version", r.Proto)

	m.next = next
	m.err = nil

	m.anubis.ServeHTTP(w, r)
	if m.err != nil {
		return m.err
	}

	return nil
}

func (m *AnubisMiddleware) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	d.Next()

	m.DefaultDifficulty = anubis.DefaultDifficulty
	m.Options.CookieExpiration = anubis.CookieDefaultExpirationTime
	m.Options.CookieSecure = true

	for nesting := d.Nesting(); d.NextBlock(nesting); {
		var err error

		switch d.Val() {
		case "difficulty":
			if !d.Next() {
				return d.ArgErr()
			}
			m.DefaultDifficulty, err = strconv.Atoi(d.Val())
			if err != nil {
				return d.WrapErr(err)
			}
		case "policy_fname":
			if !d.Next() {
				return d.ArgErr()
			}
			m.PolicyFname = d.Val()
		case "private_key":
			if !d.Next() {
				return d.ArgErr()
			}
			seed, err := hex.DecodeString(d.Val())
			if err != nil {
				return d.WrapErr(err)
			}
			m.Options.ED25519PrivateKey = ed25519.NewKeyFromSeed(seed)
		}
	} // anubis options

	if d.NextArg() {
		return d.ArgErr()
	} // too many args

	return nil
}

func parseCaddyfileHandler(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var m AnubisMiddleware
	err := m.UnmarshalCaddyfile(h.Dispenser)
	return &m, err
}
