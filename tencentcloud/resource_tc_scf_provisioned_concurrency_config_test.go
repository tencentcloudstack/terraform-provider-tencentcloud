package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
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
			{
				ResourceName:      "tencentcloud_scf_provisioned_concurrency_config.provisioned_concurrency_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccScfProvisionedConcurrencyConfig = `

resource "tencentcloud_scf_provisioned_concurrency_config" "provisioned_concurrency_config" {
  function_name = "test_function"
  qualifier = "1"
  version_provisioned_concurrency_num = 2
  namespace = "test_namespace"
  trigger_actions {
		trigger_name = ""
		trigger_provisioned_concurrency_num = 
		trigger_cron_config = ""
		provisioned_type = ""

  }
  provisioned_type = ""
  tracking_target = 
  min_capacity = 
  max_capacity = 
}

`
