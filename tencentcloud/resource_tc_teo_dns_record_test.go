package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_teo_dns_record
	resource.AddTestSweepers("tencentcloud_teo_dns_record", &resource.Sweeper{
		Name: "tencentcloud_teo_dns_record",
		F:    testSweepDnsRecord,
	})
}

func testSweepDnsRecord(region string) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cli, _ := sharedClientForRegion(region)
	client := cli.(*TencentCloudClient).apiV3Conn
	service := TeoService{client}

	zoneId := defaultZoneId

	for {
		record, err := service.DescribeTeoDnsRecord(ctx, zoneId, "")
		if err != nil {
			return err
		}

		if record == nil {
			return nil
		}

		err = service.DeleteTeoDnsRecordById(ctx, zoneId, *record.DnsRecordId)
		if err != nil {
			return err
		}
	}
}

// go test -i; go test -test.run TestAccTencentCloudTeoDnsRecord_basic -v
func TestAccTencentCloudTeoDnsRecord_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PRIVATE) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoDnsRecord,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRecordExists("tencentcloud_teo_dns_record.basic"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record.basic", "zone_id", defaultZoneId),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_dns_record.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckDnsRecordDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TeoService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_teo_dns_record" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		zoneId := idSplit[0]
		proxyId := idSplit[1]

		agents, err := service.DescribeTeoDnsRecord(ctx, zoneId, proxyId)
		if agents != nil {
			return fmt.Errorf("zone DnsRecord %s still exists", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckDnsRecordExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		zoneId := idSplit[0]
		proxyId := idSplit[1]

		service := TeoService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		agents, err := service.DescribeTeoDnsRecord(ctx, zoneId, proxyId)
		if agents == nil {
			return fmt.Errorf("zone DnsRecord %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTeoDnsRecordVar = `
variable "zone_id" {
  default = "` + defaultZoneId + `"
}

variable "zone_name" {
  default = "aaa.` + defaultZoneName + `"
}`

const testAccTeoDnsRecord = testAccTeoDnsRecordVar + `

resource "tencentcloud_teo_dns_record" "basic" {
  zone_id   = var.zone_id
  type      = "A"
  name      = var.zone_name
  content   = "150.109.8.2"
  mode      = "proxied"
  ttl       = "1"
  priority  = 1
}
`
