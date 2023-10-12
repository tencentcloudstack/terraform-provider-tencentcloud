package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslRevokeCertificateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCommon(t, ACCOUNT_TYPE_SSL)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslRevokeCertificate,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ssl_revoke_certificate_operation.revoke_certificate", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ssl_revoke_certificate_operation.revoke_certificate", "certificate_id", "8hUkH3xC"),
				),
			},
		},
	})
}

const testAccSslRevokeCertificate = `

resource "tencentcloud_ssl_revoke_certificate_operation" "revoke_certificate" {
  certificate_id = "8hUkH3xC"
}

`
