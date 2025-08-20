package dlc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudDlcUpdateStandardEngineResourceGroupResourceInformationOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccDlcUpdateStandardEngineResourceGroupResourceInformationOperation,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("tencentcloud_dlc_update_standard_engine_resource_group_resource_information_operation.example", "id"),
			),
		}},
	})
}

const testAccDlcUpdateStandardEngineResourceGroupResourceInformationOperation = `
resource "tencentcloud_dlc_update_standard_engine_resource_group_resource_information_operation" "example" {
  engine_resource_group_name = "test"
}
`
