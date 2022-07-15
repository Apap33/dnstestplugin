package dnstestplugin

import (
	"context"
	"coredns/request"
	"fmt"
	"net"
	"strings"

	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"

	"github.com/miekg/dns"
)

var log = clog.NewWithPlugin("dnstestplugin")

type DnsTestPlugin struct {
	Next plugin.Handler
}

func (d DnsTestPlugin) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	qname := state.Name()
	reply := "8.8.8.8"

	if strings.HasPrefix(state.IP(), "172.") || strings.HasPrefix(state.IP(), "127.") {
		reply = "1.1.1.1"
	}

	fmt.Printf("received query %s from %s, expected to reply %s\n", qname, state.IP(), reply)

	m := new(dns.Msg)
	m.SetReply(r)
	hdr := dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeTXT, Class: dns.ClassCHAOS, Ttl: 0}

	m.Answer = []dns.RR{&dns.A{Hdr: hdr, A: net.ParseIP(reply)}}

	w.WriteMsg(m)
	return 0, nil

}

func (d DnsTestPlugin) Name() string { return "dnstestplugin" }
