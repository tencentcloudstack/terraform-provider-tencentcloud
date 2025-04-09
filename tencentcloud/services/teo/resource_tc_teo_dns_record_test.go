package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoDnsRecordResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoDnsRecord,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_teo_dns_record.teo_dns_record", "id")),
			},
			{
				Config: testAccTeoDnsRecordUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record.teo_dns_record", "content", "1.2.3.6"),
				),
			},
			{
				Config: testAccTeoDnsRecordUpdateStatus,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record.teo_dns_record", "status", "disable"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_dns_record.teo_dns_record",
				ImportState:       true,
				ImportStateVerify: true,
			}},
	})
}

const testAccTeoDnsRecord = `
resource "tencentcloud_teo_dns_record" "teo_dns_record" {
  zone_id  = "zone-39quuimqg8r6"
  type     = "A"
  content  = "1.2.3.5"
  location = "Default"
  name     = "a.makn.cn"
  priority = 5
  ttl      = 300
  weight   = -1
  status   = "enable"
}
`

const testAccTeoDnsRecordUpdate = `
resource "tencentcloud_teo_dns_record" "teo_dns_record" {
  zone_id  = "zone-39quuimqg8r6"
  type     = "A"
  content  = "1.2.3.6"
  location = "Default"
  name     = "a.makn.cn"
  priority = 5
  ttl      = 300
  weight   = -1
  status   = "enable"
}
`

const testAccTeoDnsRecordUpdateStatus = `
resource "tencentcloud_teo_dns_record" "teo_dns_record" {
  zone_id  = "zone-39quuimqg8r6"
  type     = "A"
  content  = "1.2.3.6"
  location = "Default"
  name     = "a.makn.cn"
  priority = 5
  ttl      = 300
  weight   = -1
  status   = "disable"
}
`
