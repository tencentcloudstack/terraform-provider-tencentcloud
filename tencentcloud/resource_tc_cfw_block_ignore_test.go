package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixCfwBlockIgnoreListResource_basic -v
func TestAccTencentCloudNeedFixCfwBlockIgnoreListResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfwBlockIgnoreIp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_block_ignore_list.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_block_ignore_list.example", "direction"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_block_ignore_list.example", "comment"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_block_ignore_list.example", "start_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_block_ignore_list.example", "end_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_block_ignore_list.example", "rule_type"),
				),
			},
			{
				ResourceName:      "tencentcloud_cfw_block_ignore_list.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCfwBlockIgnoreIpUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_block_ignore_list.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_block_ignore_list.example", "direction"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_block_ignore_list.example", "comment"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_block_ignore_list.example", "start_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_block_ignore_list.example", "end_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_block_ignore_list.example", "rule_type"),
				),
			},
		},
	})
}

const testAccCfwBlockIgnoreIp = `
resource "tencentcloud_cfw_block_ignore" "example" {
  ip         = "1.1.1.1"
  direction  = 0
  comment    = "remark."
  start_time = "2023-09-01 00:00:00"
  end_time   = "2023-10-01 00:00:00"
  rule_type  = 1
}
`

const testAccCfwBlockIgnoreIpUpdate = `
resource "tencentcloud_cfw_block_ignore" "example" {
  ip         = "1.1.1.1"
  direction  = 0
  comment    = "remark update."
  start_time = "2023-09-01 00:00:00"
  end_time   = "2023-11-01 00:00:00"
  rule_type  = 1
}
`
