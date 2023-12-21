package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSslUpdateCertificateInstanceOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCommon(t, ACCOUNT_TYPE_SSL)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslUpdateCertificateInstanceOperation,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ssl_update_certificate_instance_operation.update_certificate_instance", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ssl_update_certificate_instance_operation.update_certificate_instance", "certificate_id", "AMpBxwPq"),
					resource.TestCheckResourceAttr("tencentcloud_ssl_update_certificate_instance_operation.update_certificate_instance", "old_certificate_id", "AN1Gys3l"),
					resource.TestCheckResourceAttrSet("tencentcloud_ssl_update_certificate_instance_operation.update_certificate_instance", "resource_types.0"),
				),
			},
		},
	})
}

const testAccSslUpdateCertificateInstanceOperation = `

resource "tencentcloud_ssl_update_certificate_instance_operation" "update_certificate_instance" {
  certificate_id = "AMpBxwPq"
  old_certificate_id = "AN1Gys3l"
  resource_types = ["cdn"]
}
`
