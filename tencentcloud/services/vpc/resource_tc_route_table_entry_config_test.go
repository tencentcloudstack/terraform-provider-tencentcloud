package vpc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudRouteTableEntryConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTableEntryConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_route_table_entry_config.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_route_table_entry_config.example", "route_table_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_route_table_entry_config.example", "route_item_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_route_table_entry_config.example", "disabled"),
				),
			},
			{
				Config: testAccRouteTableEntryConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_route_table_entry_config.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_route_table_entry_config.example", "route_table_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_route_table_entry_config.example", "route_item_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_route_table_entry_config.example", "disabled"),
				),
			},
			{
				ResourceName:      "tencentcloud_route_table_entry_config.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccRouteTableEntryConfig = `
resource "tencentcloud_route_table_entry_config" "example" {
  route_table_id = "rtb-8425lgjy"
  route_item_id  = "rti-4f6efqwn"
  disabled       = false
}
`

const testAccRouteTableEntryConfigUpdate = `
resource "tencentcloud_route_table_entry_config" "example" {
  route_table_id = "rtb-8425lgjy"
  route_item_id  = "rti-4f6efqwn"
  disabled       = true
}
`
