package main

import (
	"errors"
	"flag"
	"os"

	"github.com/yankeguo/nas-tools/utils"
	"github.com/yankeguo/rg"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
)

func main() {
	var err error
	defer utils.Exit(&err)
	defer rg.Guard(&err)

	var (
		optAddress   string
		optDomain    string
		optSubdomain string
	)
	flag.StringVar(&optAddress, "address", "", "ip address")
	flag.StringVar(&optDomain, "domain", "", "domain name")
	flag.StringVar(&optSubdomain, "subdomain", "", "subdomain name")
	flag.Parse()

	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = "dnspod.tencentcloudapi.com"

	client := rg.Must(dnspod.NewClient(common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	), "", clientProfile))

	var (
		recordID   uint64
		recordLine string
	)

	{
		request := dnspod.NewDescribeRecordListRequest()

		request.Domain = common.StringPtr(optDomain)
		request.Subdomain = common.StringPtr(optSubdomain)
		request.RecordType = common.StringPtr("A")

		response := rg.Must(client.DescribeRecordList(request))
		if len(response.Response.RecordList) == 0 {
			err = errors.New("no record found")
			return
		}

		recordID = *response.Response.RecordList[0].RecordId
		recordLine = *response.Response.RecordList[0].Line
	}

	{

		request := dnspod.NewModifyRecordRequest()

		request.Domain = common.StringPtr(optDomain)
		request.RecordType = common.StringPtr("A")
		request.RecordLine = common.StringPtr(recordLine)
		request.Value = common.StringPtr(optAddress)
		request.RecordId = common.Uint64Ptr(recordID)

		_ = rg.Must(client.ModifyRecord(request))
	}
}
