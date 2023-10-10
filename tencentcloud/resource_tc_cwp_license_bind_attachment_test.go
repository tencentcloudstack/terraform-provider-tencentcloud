package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixCwpLicenseBindAttachmentResource_basic -v
func TestAccTencentCloudNeedFixCwpLicenseBindAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCwpLicenseBindAttachment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cwp_license_bind_attachment.license_bind_attachment", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cwp_license_bind_attachment.license_bind_attachment", "resource_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cwp_license_bind_attachment.license_bind_attachment", "license_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cwp_license_bind_attachment.license_bind_attachment", "license_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cwp_license_bind_attachment.license_bind_attachment", "quuid"),
				),
			},
			{
				ResourceName:      "tencentcloud_cwp_license_bind_attachment.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCwpLicenseBindAttachment = `
resource "tencentcloud_cwp_license_order" "example" {
  alias        = "tf_example"
  license_type = 0
  license_num  = 1
  region_id    = 1
  project_id   = 0
  tags         = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_cwp_license_bind_attachment" "example" {
  resource_id  = tencentcloud_cwp_license_order.example.resource_id
  license_id   = tencentcloud_cwp_license_order.example.license_id
  license_type = 0
  quuid        = "2c7e5cce-1cec-4456-8d18-018f160dd987"
}
`
