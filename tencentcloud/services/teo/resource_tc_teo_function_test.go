package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoFunctionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoFunction,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function.teo_function", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.teo_function", "name", "aaa-zone-2qtuhspy7cr6-1310708577"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.teo_function", "remark", "test"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function.teo_function", "function_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.teo_function", "content", `addEventListener('fetch', e => {
  const response = new Response('Hello World!!');
  e.respondWith(response);
});
`),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_function.teo_function",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTeoFunctionUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function.teo_function", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.teo_function", "name", "aaa-zone-2qtuhspy7cr6-1310708577"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.teo_function", "remark", "test-update"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.teo_function", "content", `addEventListener('fetch', e => {
  const response = new Response('Hello World');
  e.respondWith(response);
});
`),
				),
			},
		},
	})
}

const testAccTeoFunction = `

resource "tencentcloud_teo_function" "teo_function" {
    content     = <<-EOT
        addEventListener('fetch', e => {
          const response = new Response('Hello World!!');
          e.respondWith(response);
        });
    EOT
    name        = "aaa-zone-2qtuhspy7cr6-1310708577"
    remark      = "test"
    zone_id     = "zone-2qtuhspy7cr6"
}
`
const testAccTeoFunctionUp = `

resource "tencentcloud_teo_function" "teo_function" {
    content     = <<-EOT
        addEventListener('fetch', e => {
          const response = new Response('Hello World');
          e.respondWith(response);
        });
    EOT
    name        = "aaa-zone-2qtuhspy7cr6-1310708577"
    remark      = "test-update"
    zone_id     = "zone-2qtuhspy7cr6"
}
`
// Note: Testing with an existing FunctionId requires a real pre-existing function ID.
// This scenario can be tested in a real environment with:
// 1. A pre-existing function created outside Terraform
// 2. Using that function_id in the resource configuration
// Example:
// resource "tencentcloud_teo_function" "teo_function_with_id" {
//     content     = "..."
//     name        = "..."
//     zone_id     = "..."
//     function_id = "existing-function-id-from-console"
// }
`
