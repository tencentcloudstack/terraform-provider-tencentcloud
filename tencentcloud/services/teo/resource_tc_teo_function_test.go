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

func TestAccTencentCloudTeoFunction_customFunctionId(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoFunctionWithCustomId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function.custom_function", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.custom_function", "function_id", "test-custom-function-id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.custom_function", "name", "custom-id-function"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.custom_function", "content", `addEventListener('fetch', e => {
  const response = new Response('Custom Function');
  e.respondWith(response);
});
`),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_function.custom_function",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoFunctionWithCustomId = `

resource "tencentcloud_teo_function" "custom_function" {
    function_id = "test-custom-function-id"
    content     = <<-EOT
        addEventListener('fetch', e => {
          const response = new Response('Custom Function');
          e.respondWith(response);
        });
    EOT
    name        = "custom-id-function"
    remark      = "test custom function id"
    zone_id     = "zone-2qtuhspy7cr6"
}
`

func TestAccTencentCloudTeoFunction_apiGeneratedFunctionId(t *testing.T) {
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
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function.api_generated", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function.api_generated", "function_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.api_generated", "name", "api-generated-function"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.api_generated", "content", `addEventListener('fetch', e => {
  const response = new Response('API Generated');
  e.respondWith(response);
});
`),
				),
			},
		},
	})
}

const testAccTeoFunctionWithoutId = `

resource "tencentcloud_teo_function" "api_generated" {
    content     = <<-EOT
        addEventListener('fetch', e => {
          const response = new Response('API Generated');
          e.respondWith(response);
        });
    EOT
    name        = "api-generated-function"
    remark      = "test api generated function id"
    zone_id     = "zone-2qtuhspy7cr6"
}
`
