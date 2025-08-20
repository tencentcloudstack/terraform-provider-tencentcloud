package dlc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudDlcUpdateStandardEngineResourceGroupConfigInformationOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcUpdateStandardEngineResourceGroupConfigInformationOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_update_standard_engine_resource_group_config_information_operation.example", "id"),
				),
			},
		},
	})
}

const testAccDlcUpdateStandardEngineResourceGroupConfigInformationOperation = `
resource "tencentcloud_dlc_update_standard_engine_resource_group_config_information_operation" "example" {
  engine_resource_group_name = "tf-example"
  update_conf_context {
    config_type = "StaticConfigType"
    params {
      config_item  = "spark.sql.shuffle.partitions"
      config_value = "300"
      operate      = "ADD"
    }
  }
}
`
