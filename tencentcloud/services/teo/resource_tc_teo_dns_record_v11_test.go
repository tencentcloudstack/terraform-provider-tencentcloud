package teo_test

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcteo "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudTeoDnsRecordV11Resource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTeoDnsRecordV11Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoDnsRecordV11Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeoDnsRecordV11Exists("tencentcloud_teo_dns_record_v11.dns_record"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_dns_record_v11.dns_record", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_v11.dns_record", "zone_id", "zone-xxxx"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_v11.dns_record", "name", "test"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_v11.dns_record", "type", "A"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_v11.dns_record", "content", "1.2.3.4"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_dns_record_v11.dns_record", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_dns_record_v11.dns_record", "created_on"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_dns_record_v11.dns_record", "modified_on"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_dns_record_v11.dns_record",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTeoDnsRecordV11Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeoDnsRecordV11Exists("tencentcloud_teo_dns_record_v11.dns_record"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_dns_record_v11.dns_record", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_v11.dns_record", "content", "5.6.7.8"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_v11.dns_record", "ttl", "600"),
				),
			},
		},
	})
}

func TestAccTencentCloudTeoDnsRecordV11Resource_withOptionalFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTeoDnsRecordV11Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoDnsRecordV11WithOptionalFields,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeoDnsRecordV11Exists("tencentcloud_teo_dns_record_v11.dns_record"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_dns_record_v11.dns_record", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_v11.dns_record", "name", "www"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_v11.dns_record", "type", "A"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_v11.dns_record", "content", "1.2.3.4"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_v11.dns_record", "ttl", "300"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_v11.dns_record", "weight", "50"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_v11.dns_record", "location", "Default"),
				),
			},
		},
	})
}

func TestAccTencentCloudTeoDnsRecordV11Resource_mxRecord(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTeoDnsRecordV11Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoDnsRecordV11MxRecord,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeoDnsRecordV11Exists("tencentcloud_teo_dns_record_v11.dns_record"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_dns_record_v11.dns_record", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_v11.dns_record", "type", "MX"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_v11.dns_record", "content", "mxdomain.qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_v11.dns_record", "priority", "10"),
				),
			},
		},
	})
}

func TestAccTencentCloudTeoDnsRecordV11Resource_updateImmutableFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoDnsRecordV11Basic,
			},
			{
				Config:      testAccTeoDnsRecordV11UpdateName,
				ExpectError: regexp.MustCompile("name and type fields cannot be modified after creation"),
			},
		},
	})
}

func testAccCheckTeoDnsRecordV11Destroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcteo.NewTeoService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_teo_dns_record_v11" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, "#")
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken, %s", rs.Primary.ID)
		}

		zoneId := idSplit[0]
		recordId := idSplit[1]

		record, err := service.DescribeDnsRecordById(ctx, zoneId, recordId)
		if record != nil {
			return fmt.Errorf("DNS record %s still exists", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func testAccCheckTeoDnsRecordV11Exists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		idSplit := strings.Split(rs.Primary.ID, "#")
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken, %s", rs.Primary.ID)
		}

		zoneId := idSplit[0]
		recordId := idSplit[1]

		service := svcteo.NewTeoService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		record, err := service.DescribeDnsRecordById(ctx, zoneId, recordId)
		if record == nil {
			return fmt.Errorf("DNS record %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTeoDnsRecordV11Basic = `
resource "tencentcloud_teo_dns_record_v11" "dns_record" {
    zone_id = "zone-xxxx"
    name    = "test"
    type    = "A"
    content = "1.2.3.4"
}
`

const testAccTeoDnsRecordV11Update = `
resource "tencentcloud_teo_dns_record_v11" "dns_record" {
    zone_id = "zone-xxxx"
    name    = "test"
    type    = "A"
    content = "5.6.7.8"
    ttl     = 600
}
`

const testAccTeoDnsRecordV11UpdateName = `
resource "tencentcloud_teo_dns_record_v11" "dns_record" {
    zone_id = "zone-xxxx"
    name    = "test-updated"
    type    = "A"
    content = "1.2.3.4"
}
`

const testAccTeoDnsRecordV11WithOptionalFields = `
resource "tencentcloud_teo_dns_record_v11" "dns_record" {
    zone_id  = "zone-xxxx"
    name     = "www"
    type     = "A"
    content  = "1.2.3.4"
    ttl      = 300
    weight   = 50
    location = "Default"
}
`

const testAccTeoDnsRecordV11MxRecord = `
resource "tencentcloud_teo_dns_record_v11" "dns_record" {
    zone_id  = "zone-xxxx"
    name     = "@"
    type     = "MX"
    content  = "mxdomain.qq.com"
    priority = 10
}
`
