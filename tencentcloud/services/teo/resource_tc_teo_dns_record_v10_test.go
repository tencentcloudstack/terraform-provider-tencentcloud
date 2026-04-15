package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoDnsRecordV10_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoDnsRecordV10,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_teo_dns_record_v10.teo_dns_record_v10", "id")),
			},
			{
				Config: testAccTeoDnsRecordV10Update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_v10.teo_dns_record_v10", "content", "1.2.3.6"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_dns_record_v10.teo_dns_record_v10",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudTeoDnsRecordV10_update(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoDnsRecordV10,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_teo_dns_record_v10.teo_dns_record_v10", "id")),
			},
			{
				Config: testAccTeoDnsRecordV10Update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_v10.teo_dns_record_v10", "content", "1.2.3.6"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_v10.teo_dns_record_v10", "ttl", "600"),
				),
			},
		},
	})
}

func TestAccTencentCloudTeoDnsRecordV10_import(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoDnsRecordV10,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_teo_dns_record_v10.teo_dns_record_v10", "id")),
			},
			{
				ResourceName:      "tencentcloud_teo_dns_record_v10.teo_dns_record_v10",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudTeoDnsRecordV10_delete(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoDnsRecordV10,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_teo_dns_record_v10.teo_dns_record_v10", "id")),
			},
			{
				Config:  "",
				Destroy: true,
				Check:   resource.ComposeTestCheckFunc(),
			},
		},
	})
}

const testAccTeoDnsRecordV10 = `
resource "tencentcloud_teo_dns_record_v10" "teo_dns_record_v10" {
  zone_id  = "zone-39quuimqg8r6"
  type     = "A"
  content  = "1.2.3.5"
  location = "Default"
  name     = "a.makn.cn"
  priority = 5
  ttl      = 300
  weight   = -1
}
`

const testAccTeoDnsRecordV10Update = `
resource "tencentcloud_teo_dns_record_v10" "teo_dns_record_v10" {
  zone_id  = "zone-39quuimqg8r6"
  type     = "A"
  content  = "1.2.3.6"
  location = "Default"
  name     = "a.makn.cn"
  priority = 5
  ttl      = 600
  weight   = -1
}
`
