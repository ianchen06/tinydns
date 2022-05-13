package main

import (
	"fmt"
	"log"
	"net"

	"github.com/miekg/dns"
)

var domainsToAddresses map[string]string = map[string]string{
	"google.com.": "1.2.3.4",
}

type handler struct{}

func (h *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := dns.Msg{}
	msg.SetReply(r)

	switch r.Question[0].Qtype {
	case dns.TypeA:
		msg.Authoritative = true
		domain := msg.Question[0].Name
		fmt.Println(domain)
		addr, ok := domainsToAddresses[domain]
		if ok {
			msg.Answer = append(msg.Answer, &dns.A{Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				A: net.ParseIP(addr)})
		}
	}
	w.WriteMsg(&msg)
}

func main() {
	srv := &dns.Server{Addr: ":53", Net: "udp"}
	srv.Handler = &handler{}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to set udp listener %s\n", err.Error())
	}
}
