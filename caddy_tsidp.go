package caddy_tsidp

import (
	// "crypto/ed25519"
	// "encoding/hex"
	// "fmt"
	"net"
	"net/http"
	// "strconv"

	srvtsidp "github.com/tailscale/tsidp/server"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
	"tailscale.com/client/local"
	// "tailscale.com/ipn/ipnstate"
)

func init() {
	caddy.RegisterModule(TsidpMiddleware{})
	httpcaddyfile.RegisterHandlerDirective("tsidp", parseCaddyfileHandler)
	httpcaddyfile.RegisterDirectiveOrder("tsidp", httpcaddyfile.Before, "push")
}

func (TsidpMiddleware) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.tsidp",
		New: func() caddy.Module { return new(TsidpMiddleware) },
	}
}

type TsidpMiddleware struct {
	PolicyFname       string            `json:"policy_fname,omitempty"`
	DefaultDifficulty int               `json:"default_difficulty,omitempty"`

	// Options
	StateDir          string            `json:"state_dir"`
	EnableSTS         bool              `json:"enable_sts,omitempty"`
	CanonicalHostname string            `json:"canonical_hostname,omitempty"`

	lc     *local.Client
	tsidp  *srvtsidp.IDPServer
	log    *zap.Logger
	next   caddyhttp.Handler
	err    error
}

// Interface guards
var (
	_ caddyhttp.MiddlewareHandler = (*TsidpMiddleware)(nil)
	_ caddyfile.Unmarshaler       = (*TsidpMiddleware)(nil)
	_ caddy.Provisioner           = (*TsidpMiddleware)(nil)
)

func (m *TsidpMiddleware) Provision(ctx caddy.Context) error {
	var (
		// st          *ipnstate.Status
		// err         error
	)
	m.log = ctx.Logger()
	// m.Options.Logger = ctx.Slogger()

	// m.log.Debug("loading tsidp policies", zap.String("policy_file", m.PolicyFname), zap.Int("default_difficulty", m.DefaultDifficulty))
	// policy, err := libtsidp.LoadPoliciesOrDefault(ctx, m.PolicyFname, m.DefaultDifficulty, ctx.Logger().Level().String(), false)
	// if err != nil {
	// 	return fmt.Errorf("failed to load tsidp policies from '%s': %w", m.PolicyFname, err)
	// }

	// m.Options.Next = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	if m.err = m.next.ServeHTTP(w, r); err != nil {
	// 		m.log.Debug("received error from next handler", zap.Error(err))
	// 	}
	// })
	// m.Options.Policy = policy
	// m.tsidp, err = libtsidp.New(m.Options)
	m.lc = &local.Client{}
	// st, err = m.lc.StatusWithoutPeers(ctx)
	// if err != nil {
	// 	return err
	// }

	m.tsidp = srvtsidp.New(
		m.lc,
		m.StateDir,
		true, // funnel
		true, // localTSMode
		m.EnableSTS,
	)

	if m.CanonicalHostname != "" {
		// TODO: Split hostname
		m.tsidp.SetServerURL(m.CanonicalHostname, 443)
	}

	if err := m.tsidp.LoadFunnelClients(); err != nil {
		return err
	}

	return nil
}

func (m *TsidpMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	remoteHost, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return err
	}
	r.Header.Set("X-Real-Ip", remoteHost)
	r.Header.Set("X-Http-Version", r.Proto)
	r.Header.Set("X-Forwarded-For", r.RemoteAddr)

	// FIXME: Call .SetServerURL

	m.tsidp.ServeHTTP(w, r)

	return nil
}

func (m *TsidpMiddleware) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	d.Next()

	for nesting := d.Nesting(); d.NextBlock(nesting); {

		switch d.Val() {
		case "state_dir":
			if !d.Next() {
				return d.ArgErr()
			}
			m.StateDir = d.Val()
		case "canonical_hostname":
			if !d.Next() {
				return d.ArgErr()
			}
			m.CanonicalHostname = d.Val()
		case "enable_sts":
			if d.NextArg() {
				return d.ArgErr()
			}
			m.EnableSTS = true
		}
	} // tsidp options

	if d.NextArg() {
		return d.ArgErr()
	} // too many args

	return nil
}

func parseCaddyfileHandler(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var m TsidpMiddleware
	err := m.UnmarshalCaddyfile(h.Dispenser)
	return &m, err
}
