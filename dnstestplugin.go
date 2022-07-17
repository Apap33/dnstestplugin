package dnstestplugin

import (
	"context"
	"fmt"
	"net"

	"github.com/coredns/coredns/request"

	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"

	"github.com/go-redis/redis"
	"github.com/miekg/dns"
)

var log = clog.NewWithPlugin("dnstestplugin")

type DnsTestPlugin struct {
	Next plugin.Handler
}

func (d DnsTestPlugin) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	// redisctx := context.Background()

	state := request.Request{W: w, Req: r}
	qname := state.Name()
	m := new(dns.Msg)
	m.SetReply(r)
	fmt.Println("Query name: ", qname)
	hdr := dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60}

	val, err := rdb.Get("testavi").Result()
	if err != nil {
		fmt.Errorf("Redis error %s", err)
		m.SetRcode(state.Req, dns.RcodeNameError)
		state.SizeAndDo(m)
		_ = state.W.WriteMsg(m)
		return dns.RcodeSuccess, err
	}
	m.Answer = []dns.RR{&dns.A{Hdr: hdr, A: net.ParseIP(val)}}

	w.WriteMsg(m)
	return 0, nil

}

func (d DnsTestPlugin) Name() string { return "dnstestplugin" }
