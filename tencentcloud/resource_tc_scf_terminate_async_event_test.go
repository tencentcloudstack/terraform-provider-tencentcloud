package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixScfTerminateAsyncEventResource_basic(t *testing.T) {
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
		},
	})
}

const testAccScfTerminateAsyncEvent = `

resource "tencentcloud_scf_terminate_async_event" "terminate_async_event" {
  function_name = "keep-1676351130"
  invoke_request_id = "9de9405a-e33a-498d-bb59-e80b7bed1191"
  namespace     = "default"
  grace_shutdown = true
}

`
