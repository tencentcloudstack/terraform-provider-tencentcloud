package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudScfFunctionEventInvokeConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccScfFunctionEventInvokeConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_scf_function_event_invoke_config.function_event_invoke_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_scf_function_event_invoke_config.function_event_invoke_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccScfFunctionEventInvokeConfig = `

resource "tencentcloud_scf_function_event_invoke_config" "function_event_invoke_config" {
  async_trigger_config {
		retry_config {
			retry_num = 2
		}
		msg_t_t_l = 24

  }
  function_name = "test_function"
  namespace = "test_namespace"
}

`
