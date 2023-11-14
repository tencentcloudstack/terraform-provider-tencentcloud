package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsExecuteFunctionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsExecuteFunction,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mps_execute_function.execute_function", "id")),
			},
			{
				ResourceName:      "tencentcloud_mps_execute_function.execute_function",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsExecuteFunction = `

resource "tencentcloud_mps_execute_function" "execute_function" {
  function_name = ""
  function_arg = ""
}

`
