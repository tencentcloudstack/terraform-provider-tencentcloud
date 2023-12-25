package ssl_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslRevokeCertificateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_SSL)
		},
		Providers: tcacctest.AccProviders,
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
