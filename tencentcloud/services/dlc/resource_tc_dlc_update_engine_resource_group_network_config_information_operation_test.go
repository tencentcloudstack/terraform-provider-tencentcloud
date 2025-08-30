package dlc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudDlcUpdateEngineResourceGroupNetworkConfigInformationOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcUpdateEngineResourceGroupNetworkConfigInformationOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_update_engine_resource_group_network_config_information_operation.example", "id"),
				),
			},
		},
	})
}

const testAccDlcUpdateEngineResourceGroupNetworkConfigInformationOperation = `
resource "tencentcloud_dlc_update_engine_resource_group_network_config_information_operation" "example" {
  engine_resource_group_id = "rg-b6fxxxxxx2a0"
}
`
