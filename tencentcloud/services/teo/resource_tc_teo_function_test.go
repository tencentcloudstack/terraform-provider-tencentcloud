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

func TestAccTencentCloudTeoFunctionResource_WithFunctionId(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoFunctionWithId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function.teo_function_with_id", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.teo_function_with_id", "function_id", "test-function-id-123"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.teo_function_with_id", "name", "test-with-id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.teo_function_with_id", "remark", "test-with-id-remark"),
				),
			},
		},
	})
}

const testAccTeoFunctionWithId = `

resource "tencentcloud_teo_function" "teo_function_with_id" {
    content     = <<-EOT
        addEventListener('fetch', e => {
          const response = new Response('Hello World with ID');
          e.respondWith(response);
        });
    EOT
    name        = "test-with-id"
    remark      = "test-with-id-remark"
    zone_id     = "zone-2qtuhspy7cr6"
    function_id = "test-function-id-123"
}
`

func TestAccTencentCloudTeoFunctionResource_WithoutFunctionId(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoFunctionWithoutId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function.teo_function_without_id", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function.teo_function_without_id", "function_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.teo_function_without_id", "name", "test-without-id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.teo_function_without_id", "remark", "test-without-id-remark"),
				),
			},
		},
	})
}

const testAccTeoFunctionWithoutId = `

resource "tencentcloud_teo_function" "teo_function_without_id" {
    content     = <<-EOT
        addEventListener('fetch', e => {
          const response = new Response('Hello World without ID');
          e.respondWith(response);
        });
    EOT
    name        = "test-without-id"
    remark      = "test-without-id-remark"
    zone_id     = "zone-2qtuhspy7cr6"
}
`

