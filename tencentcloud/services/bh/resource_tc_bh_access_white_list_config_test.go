package bh_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudBhAccessWhiteListConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBhAccessWhiteListConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_bh_access_white_list_config.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_bh_access_white_list_config.example", "allow_any"),
					resource.TestCheckResourceAttrSet("tencentcloud_bh_access_white_list_config.example", "allow_auto"),
				),
			},
			{
				Config: testAccBhAccessWhiteListConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_bh_access_white_list_config.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_bh_access_white_list_config.example", "allow_any"),
					resource.TestCheckResourceAttrSet("tencentcloud_bh_access_white_list_config.example", "allow_auto"),
				),
			},
			{
				ResourceName:      "tencentcloud_bh_access_white_list_config.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccBhAccessWhiteListConfig = `
resource "tencentcloud_bh_access_white_list_config" "example" {
  allow_any  = false
  allow_auto = false
}
`

const testAccBhAccessWhiteListConfigUpdate = `
resource "tencentcloud_bh_access_white_list_config" "example" {
  allow_any  = true
  allow_auto = true
}
`
