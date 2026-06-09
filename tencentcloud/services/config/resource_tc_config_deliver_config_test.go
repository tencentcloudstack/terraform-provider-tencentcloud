package config_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudConfigDeliverConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigDeliverConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_config_deliver_config.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_config_deliver_config.example", "status", "1"),
					resource.TestCheckResourceAttr("tencentcloud_config_deliver_config.example", "deliver_type", "COS"),
					resource.TestCheckResourceAttr("tencentcloud_config_deliver_config.example", "deliver_content_type", "3"),
				),
			},
			{
				Config: testAccConfigDeliverConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_config_deliver_config.example", "status", "0"),
				),
			},
		},
	})
}

// NOTE: Replace target_arn with a real COS/CLS ARN in your environment.
const testAccConfigDeliverConfig = `
resource "tencentcloud_config_deliver_config" "example" {
  status               = 1
  deliver_name         = "tf-example-deliver"
  target_arn           = "qcs::cos:ap-guangzhou:uin/100000005287:prefix/1307050748/my-config-bucket"
  deliver_prefix       = "config"
  deliver_type         = "COS"
  deliver_content_type = 3
}
`

const testAccConfigDeliverConfigUpdate = `
resource "tencentcloud_config_deliver_config" "example" {
  status = 0
}
`
