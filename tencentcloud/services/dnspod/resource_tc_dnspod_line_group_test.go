package dnspod_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDnspodLineGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDnspodLineGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dnspod_line_group.line_group", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_line_group.line_group", "name", "test_group"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_line_group.line_group", "lines.#", "2"),
				),
			},
			{
				ResourceName:      "tencentcloud_dnspod_line_group.line_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDnspodLineGroupUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dnspod_line_group.line_group", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_line_group.line_group", "name", "test_group_updated"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_line_group.line_group", "lines.#", "3"),
				),
			},
		},
	})
}

const testAccDnspodLineGroup = `
resource "tencentcloud_dnspod_line_group" "line_group" {
  domain = "iac-tf.cloud"
  name   = "test_group"
  lines  = ["电信", "移动"]
}
`

const testAccDnspodLineGroupUpdate = `
resource "tencentcloud_dnspod_line_group" "line_group" {
  domain = "iac-tf.cloud"
  name   = "test_group_updated"
  lines  = ["电信", "移动", "联通"]
}
`
