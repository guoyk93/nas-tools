package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/cloudflare/cloudflare-go"
	"github.com/yankeguo/nas-tools/utils"
	"github.com/yankeguo/rg"
)

func main() {
	var err error
	defer utils.Exit(&err)
	defer rg.Guard(&err)

	var (
		optAddress string
		optDomain  string
	)
	flag.StringVar(&optAddress, "address", "", "ip address")
	flag.StringVar(&optDomain, "domain", "", "domain name")
	flag.Parse()

	optAddress = strings.TrimSpace(optAddress)
	optDomain = strings.TrimSpace(optDomain)

	if optAddress == "" || optDomain == "" {
		flag.Usage()
		return
	}

	cf := rg.Must(cloudflare.New(os.Getenv("CLOUDFLARE_API_KEY"), os.Getenv("CLOUDFLARE_API_EMAIL")))

	rc := cloudflare.ZoneIdentifier(os.Getenv("CLOUDFLARE_ZONE_ID"))

	ctx := context.Background()

	records, info := rg.Must2(cf.ListDNSRecords(ctx, rc, cloudflare.ListDNSRecordsParams{
		Name: optDomain,
	}))

	log.Print(records, info)
}
