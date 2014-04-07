package main

import (
	"github.com/hashicorp/serf/client"
	"github.com/miekg/dns"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"fmt"
)

func MakeRR(s string) dns.RR { r, _ := dns.NewRR(s); return r }

const SOA string = "@ SOA ns1.example.fake. hostmaster.example.fake. 2002040800 1800 900 0604800 604800"

func MakeSOA() {
	z := "serf."
	rr := MakeRR("$ORIGIN serf.\n" + SOA)
	rrx := rr.(*dns.SOA) // Needed to create the actual RR, and not an reference.
	fmt.Println("Making zone for " + z)
	fmt.Println("With a RR of ")
	fmt.Print(rr)
	fmt.Println()
	dns.HandleFunc(z, func(w dns.ResponseWriter, r *dns.Msg) {
		m := new(dns.Msg)
		m.SetReply(r)
		m.Authoritative = true
		m.Ns = []dns.RR{rrx}
		w.WriteMsg(m)
		fmt.Println(r)
	})
}

func MakeSerfRecords() {
	serf_client, err := client.NewRPCClient("127.0.0.1:7373")
	fmt.Println(err)
	fmt.Println(serf_client)
	members, _ := serf_client.Members()
	fmt.Println(members)
	for _, member := range members {

		// Every serf member gets its own A record
		host := member.Name
		ip   := member.Addr.String()
		record := host +  ".serf. IN A " + ip + "\n"
		rr := MakeRR(record)
		rrx := rr.(*dns.A) // Needed to create the actual RR, and not an reference.
		dns.HandleFunc(host + ".serf.", func(w dns.ResponseWriter, r *dns.Msg) {
			m := new(dns.Msg)
			m.SetReply(r)
			m.Authoritative = true
			//m.Ns = []dns.RR{rrx}
			m.Answer = []dns.RR{rrx}
			w.WriteMsg(m)
			fmt.Println(r)
		})
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 4)
	MakeSOA()
	MakeSerfRecords()
	go func() {
		err := dns.ListenAndServe(":8053", "tcp", nil)
		if err != nil {
			log.Fatal("Failed to set tcp listener %s\n", err.Error())
		}
	}()
	go func() {
		err := dns.ListenAndServe(":8053", "udp", nil)
		if err != nil {
			log.Fatal("Failed to set udp listener %s\n", err.Error())
		}
	}()
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case s := <-sig:
			log.Fatalf("Signal (%d) received, stopping\n", s)
		}
	}
}
