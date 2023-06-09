package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudScfProvisionedConcurrencyConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccScfProvisionedConcurrencyConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_scf_provisioned_concurrency_config.provisioned_concurrency_config", "id")),
			},
		},
	})
}

const testAccScfProvisionedConcurrencyConfig = `

resource "tencentcloud_scf_provisioned_concurrency_config" "provisioned_concurrency_config" {
  function_name                       = "keep-1676351130"
  qualifier                           = "2"
  version_provisioned_concurrency_num = 2
  namespace                           = "default"
  trigger_actions {
    trigger_name                        = "test"
    trigger_provisioned_concurrency_num = 2
    trigger_cron_config                 = "29 45 12 29 05 * 2023"
    provisioned_type                    = "Default"
  }
  provisioned_type                    = "Default"
  tracking_target                     = 0.5
  min_capacity                        = 1
  max_capacity                        = 2
}

`
