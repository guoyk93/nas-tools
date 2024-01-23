package main

import (
	"context"
	"errors"
	"flag"
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

	cf := rg.Must(cloudflare.NewWithAPIToken(os.Getenv("CLOUDFLARE_API_TOKEN")))

	rc := cloudflare.ZoneIdentifier(os.Getenv("CLOUDFLARE_ZONE_ID"))

	ctx := context.Background()

	records, _ := rg.Must2(cf.ListDNSRecords(ctx, rc, cloudflare.ListDNSRecordsParams{
		Name: optDomain,
	}))

	var record cloudflare.DNSRecord

	for _, _record := range records {
		if _record.Type == "A" {
			record = _record
			break
		}
	}

	if record.ID == "" {
		err = errors.New("no A record found")
		return
	}

	if record.Content == optAddress {
		return
	}

	rg.Must(cf.UpdateDNSRecord(ctx, rc, cloudflare.UpdateDNSRecordParams{
		ID:      record.ID,
		Name:    record.Name,
		Type:    "A",
		Content: optAddress,
		Proxied: utils.Ptr(false),
		TTL:     300,
	}))
}
