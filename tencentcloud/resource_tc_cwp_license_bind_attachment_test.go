package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudCwpLicenseBindAttachmentResource_basic -v
func TestAccTencentCloudCwpLicenseBindAttachmentResource_basic(t *testing.T) {
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
resource "tencentcloud_cwp_license_bind_attachment" "example" {
  resource_id  = ""
  license_id   = 0
  license_type = 0
  quuid        = ""
}
`
