package bh_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudBhAccessWhiteListRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBhAccessWhiteListRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_bh_access_white_list_rule.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_bh_access_white_list_rule.example", "source"),
					resource.TestCheckResourceAttrSet("tencentcloud_bh_access_white_list_rule.example", "remark"),
				),
			},
			{
				Config: testAccBhAccessWhiteListRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_bh_access_white_list_rule.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_bh_access_white_list_rule.example", "source"),
					resource.TestCheckResourceAttrSet("tencentcloud_bh_access_white_list_rule.example", "remark"),
				),
			},
			{
				ResourceName:      "tencentcloud_bh_access_white_list_rule.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccBhAccessWhiteListRule = `
resource "tencentcloud_bh_access_white_list_rule" "example" {
  source = "1.1.1.1"
  remark = "remark."
}
`

const testAccBhAccessWhiteListRuleUpdate = `
resource "tencentcloud_bh_access_white_list_rule" "example" {
  source = "2.2.2.2"
  remark = "remark update."
}
`
