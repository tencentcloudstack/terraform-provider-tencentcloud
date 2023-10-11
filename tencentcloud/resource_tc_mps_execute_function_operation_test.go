package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMpsExecuteFunctionOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsExecuteFunctionOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mps_execute_function_operation.execute_function_operation", "id")),
			},
			{
				ResourceName:      "tencentcloud_mps_execute_function_operation.execute_function_operation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsExecuteFunctionOperation = `

resource "tencentcloud_mps_execute_function_operation" "execute_function_operation" {
  function_name = ""
  function_arg = ""
}

`
