package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudScfInvokeFunctionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccScfInvokeFunction,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_scf_invoke_function.invoke_function", "id")),
			},
			{
				ResourceName:      "tencentcloud_scf_invoke_function.invoke_function",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccScfInvokeFunction = `

resource "tencentcloud_scf_invoke_function" "invoke_function" {
  function_name = "test_function"
  invocation_type = ""
  qualifier = ""
  client_context = ""
  log_type = ""
  namespace = "test_namespace"
  routing_key = ""
}

`
