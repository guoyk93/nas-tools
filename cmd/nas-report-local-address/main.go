package main

import (
	"flag"
	"log"
	"net"

	"github.com/yankeguo/nas-tools/utils"
	"github.com/yankeguo/rg"
)

func main() {
	var err error
	defer utils.Exit(&err)
	defer rg.Guard(&err)

	var (
		optInterface string
	)
	flag.StringVar(&optInterface, "if", "ovs_eth0", "interface name")
	flag.Parse()

	iface := rg.Must(net.InterfaceByName(optInterface))

	for _, addr := range rg.Must(iface.Addrs()) {
		log.Println(addr.String())
	}
}
