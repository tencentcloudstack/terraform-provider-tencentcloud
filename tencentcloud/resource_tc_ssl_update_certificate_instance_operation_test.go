package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslUpdateCertificateInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCommon(t, ACCOUNT_TYPE_SSL)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslUpdateCertificateInstance,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ssl_update_certificate_instance_operation.update_certificate_instance", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ssl_update_certificate_instance_operation.update_certificate_instance", "certificate_id", "7REHyHM1"),
					resource.TestCheckResourceAttr("tencentcloud_ssl_update_certificate_instance_operation.update_certificate_instance", "old_certificate_id", "9D3qRt7r"),
					resource.TestCheckResourceAttr("tencentcloud_ssl_update_certificate_instance_operation.update_certificate_instance", "resource_types.0", "cdn"),
				),
			},
		},
	})
}

const testAccSslUpdateCertificateInstance = `

resource "tencentcloud_ssl_update_certificate_instance_operation" "update_certificate_instance" {
  certificate_id = "7REHyHM1"
  old_certificate_id = "9D3qRt7r"
  resource_types = ["cdn"]
}
`
