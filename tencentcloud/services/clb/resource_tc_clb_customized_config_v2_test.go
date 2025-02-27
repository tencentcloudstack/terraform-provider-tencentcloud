package clb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudClbCustomizedConfigV2_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClbLogsetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbCustomizedConfigV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbLogsetExists("tencentcloud_clb_customized_config_v2.example"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_customized_config_v2.example", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_clb_customized_config_v2.example", "name", "clb_custom_config"),
				),
			},
			{
				ResourceName:      "tencentcloud_clb_customized_config_v2.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccClbCustomizedConfigV2_basic = `
resource "tencentcloud_clb_log_set" "test_logset" {
}
`
