package teo_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoFunctionV2Resource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoFunctionV2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function_v2.teo_function_v2", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_v2.teo_function_v2", "name", "aaa-zone-2qtuhspy7cr6-1310708577"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_v2.teo_function_v2", "remark", "test"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_v2.teo_function_v2", "content", `addEventListener('fetch', e => {
  const response = new Response('Hello World!!');
  e.respondWith(response);
});
`),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_function_v2.teo_function_v2",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTeoFunctionV2Up,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function_v2.teo_function_v2", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_v2.teo_function_v2", "name", "aaa-zone-2qtuhspy7cr6-1310708577"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_v2.teo_function_v2", "remark", "test-update"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_v2.teo_function_v2", "content", `addEventListener('fetch', e => {
  const response = new Response('Hello World');
  e.respondWith(response);
});
`),
				),
			},
		},
	})
}

func TestAccTencentCloudTeoFunctionV2Resource_validation(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccTeoFunctionV2InvalidName,
				ExpectError: regexp.MustCompile(".*invalid.*"),
			},
		},
	})
}

func TestAccTencentCloudTeoFunctionV2Resource_updateImmutableName(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoFunctionV2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function_v2.teo_function_v2", "id"),
				),
			},
			{
				Config:      testAccTeoFunctionV2UpdateName,
				ExpectError: regexp.MustCompile(".*cannot be changed.*"),
			},
		},
	})
}

const testAccTeoFunctionV2 = `

resource "tencentcloud_teo_function_v2" "teo_function_v2" {
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
const testAccTeoFunctionV2Up = `

resource "tencentcloud_teo_function_v2" "teo_function_v2" {
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

const testAccTeoFunctionV2InvalidName = `

resource "tencentcloud_teo_function_v2" "teo_function_v2" {
    content     = <<-EOT
        addEventListener('fetch', e => {
          const response = new Response('Hello World!!');
          e.respondWith(response);
        });
    EOT
    name        = "AAA-ZONE-INVALID"
    remark      = "test"
    zone_id     = "zone-2qtuhspy7cr6"
}
`

const testAccTeoFunctionV2UpdateName = `

resource "tencentcloud_teo_function_v2" "teo_function_v2" {
    content     = <<-EOT
        addEventListener('fetch', e => {
          const response = new Response('Hello World!!');
          e.respondWith(response);
        });
    EOT
    name        = "bbb-zone-2qtuhspy7cr6-1310708577"
    remark      = "test"
    zone_id     = "zone-2qtuhspy7cr6"
}
`
