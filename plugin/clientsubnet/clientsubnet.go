package clientsubnet

import (
	"context"
	"net"

	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

const name = "clientsubnet"

// Whoami is a plugin that returns your IP address, port and the protocol used for connecting
// to CoreDNS.
type ClientSubnet struct {
	clientIP string
	Next     plugin.Handler
}

// ServeDNS implements the plugin.Handler interface.
func (wh ClientSubnet) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	o := setupEdns0Opt(r)
	o.Option = append(o.Option, &dns.EDNS0_SUBNET{
		Code:          dns.EDNS0SUBNET,
		Family:        1,
		SourceNetmask: 32,
		SourceScope:   0,
		Address:       net.ParseIP(wh.clientIP).To4(),
	})

	return plugin.NextOrFailure(wh.Name(), wh.Next, ctx, w, r)
}

// Name implements the Handler interface.
func (wh ClientSubnet) Name() string { return name }

// setupEdns0Opt will retrieve the EDNS0 OPT or create it if it does not exist.
func setupEdns0Opt(r *dns.Msg) *dns.OPT {
	o := r.IsEdns0()
	if o == nil {
		r.SetEdns0(4096, false)
		o = r.IsEdns0()
	}
	return o
}
