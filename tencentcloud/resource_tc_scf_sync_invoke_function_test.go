package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudScfSyncInvokeFunctionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccScfSyncInvokeFunction,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_scf_sync_invoke_function.sync_invoke_function", "id")),
			},
			{
				ResourceName:      "tencentcloud_scf_sync_invoke_function.sync_invoke_function",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccScfSyncInvokeFunction = `

resource "tencentcloud_scf_sync_invoke_function" "sync_invoke_function" {
  function_name = "test_function"
  qualifier = ""
  event = ""
  log_type = ""
  namespace = ""
  routing_key = ""
}

`
