package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDnspodCustom_lineResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDnspodCustom_line,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dnspod_custom_line.custom_line", "id")),
			},
			{
				ResourceName:      "tencentcloud_dnspod_custom_line.custom_line",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDnspodCustom_line = `

resource "tencentcloud_dnspod_custom_line" "custom_line" {
  domain = "dnspod.com"
  name = "testline8"
  area = "6.6.6.1-6.6.6.3"
  domain_id = 1009
}

`
