package main

import (
	"flag"
	"fmt"
	"net"
	"strings"
)

var (
	commentsEnabled = flag.Bool("comments", false, "Enable comments after post")

	commentsLimit = flag.Int("limit", 10, "Comments number per page")

	commentsServices = &AddrList{}
)

// go build flag.go
// ./flag -comments=true -servers="127.0.0.1:8081,127.0.0.1:8082"
// ./flag --help

func init() {
	flag.Var(commentsServices, "servers", "Addresses")
}

func main() {
	flag.Parse()

	// fmt.Printf("%v\n%v\n%v\n", *commentsEnabled, *commentsLimit, *commentsServices)
	if *commentsEnabled {
		fmt.Println("Comments per page", *commentsLimit)
		fmt.Println("Comments services", *commentsServices)
	} else {
		fmt.Println("Comments disabled")
	}
}

type AddrList []string

var _ flag.Value = &AddrList{}

func (v *AddrList) String() string {
	return fmt.Sprint(*v)
}

// in 127.0.0.1:8081,127.0.0.1:8082
// v
func (v *AddrList) Set(in string) error {
	for _, addr := range strings.Split(in, ",") {
		ipRaw, _, err := net.SplitHostPort(addr)
		if err != nil {
			return fmt.Errorf("bad addr %v", addr)
		}
		ip := net.ParseIP(ipRaw)
		if ip.To4() == nil {
			return fmt.Errorf("invalid ipv4 addr %v", addr)
		}
		*v = append(*v, addr)
	}
	return nil
}
