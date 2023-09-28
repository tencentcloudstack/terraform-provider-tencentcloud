package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSslUpdateCertificateInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslUpdateCertificateInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ssl_update_certificate_instance.update_certificate_instance", "id")),
			},
		},
	})
}

const testAccSslUpdateCertificateInstance = `

resource "tencentcloud_ssl_update_certificate_instance" "update_certificate_instance" {
  certificate_id = "8x1eUSSl"
  old_certificate_id = "8xNdi2ig"
  resource_types = ["cdn"]
}
`
