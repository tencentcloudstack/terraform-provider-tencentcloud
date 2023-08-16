package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSslRevokeCertificateOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslRevokeCertificateOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ssl_revoke_certificate_operation.revoke_certificate_operation", "id")),
			},
			{
				ResourceName:      "tencentcloud_ssl_revoke_certificate_operation.revoke_certificate_operation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSslRevokeCertificateOperation = `

resource "tencentcloud_ssl_revoke_certificate_operation" "revoke_certificate_operation" {
  certificate_id = "7Me1pCxd"
  reason = "xx"
}

`
