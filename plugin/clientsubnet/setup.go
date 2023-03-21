package clientsubnet

import (
	"net"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/log"
)

func init() { plugin.Register("clientsubnet", setup) }

func setup(c *caddy.Controller) error {
	if c.Next() {
		clientIp := ""
		args := c.RemainingArgs()
		if len(args) > 0 {
			if net.ParseIP(args[0]) != nil {
				clientIp = args[0]
			}
		} else {
			log.Warning("invalid clientsubnet config")
		}

		if len(clientIp) > 0 {
			dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
				return ClientSubnet{
					clientIP: clientIp,
					Next:     next,
				}
			})
		}
	}

	return nil
}
