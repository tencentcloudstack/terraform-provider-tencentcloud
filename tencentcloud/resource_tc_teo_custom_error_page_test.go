package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTeoCustomErrorPage_basic -v
func TestAccTencentCloudTeoCustomErrorPage_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PRIVATE) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoCustomErrorPage,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_custom_error_page.basic", "id"),
				),
			},
		},
	})
}

const testAccTeoCustomErrorPageVar = `
variable "zone_id" {
  default = "` + defaultZoneId + `"
}

variable "zone_name" {
  default = "` + defaultZoneName + `"
}`

const testAccTeoCustomErrorPage = testAccTeoCustomErrorPageVar + `

resource "tencentcloud_teo_custom_error_page" "basic" {
  zone_id = var.zone_id
  entity  = var.zone_name

  name    = "test"
  content = "<html lang='en'><body><div><p>test content</p></div></body></html>"
}

`
