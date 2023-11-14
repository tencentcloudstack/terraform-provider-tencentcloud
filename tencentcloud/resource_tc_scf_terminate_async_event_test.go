package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudScfTerminateAsyncEventResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccScfTerminateAsyncEvent,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_scf_terminate_async_event.terminate_async_event", "id")),
			},
			{
				ResourceName:      "tencentcloud_scf_terminate_async_event.terminate_async_event",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccScfTerminateAsyncEvent = `

resource "tencentcloud_scf_terminate_async_event" "terminate_async_event" {
  function_name = "test"
  invoke_request_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  namespace = "testNamespace"
  grace_shutdown = true
}

`
