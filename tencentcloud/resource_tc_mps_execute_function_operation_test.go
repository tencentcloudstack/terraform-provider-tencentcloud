package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// NeedFix: reserved for the custom support
func TestAccTencentCloudNeedFixMpsExecuteFunctionOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsExecuteFunctionOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_execute_function_operation.operation", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mps_execute_function_operation.operation", "function_name", "ExampleFunc"),
					resource.TestCheckResourceAttr("tencentcloud_mps_execute_function_operation.operation", "function_arg", "arg1"),
				),
			},
		},
	})
}

const testAccMpsExecuteFunctionOperation = `

resource "tencentcloud_mps_execute_function_operation" "operation" {
  function_name = "ExampleFunc"
  function_arg = "arg1"
}

`
