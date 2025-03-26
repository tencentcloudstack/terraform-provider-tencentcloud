package clb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudClbCustomizedConfigV2_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClbCustomizedConfigV2_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_clb_customized_config_v2.example", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_clb_customized_config_v2.example", "config_name", "tf-example"),
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
resource "tencentcloud_clb_customized_config_v2" "example" {
  config_content = "client_max_body_size 224M;\r\nclient_body_timeout 60s;"
  config_name    = "tf-example"
  config_type    = "LOCATION"
}
`
