package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslUpdateCertificateRecordRetryResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCommon(t, ACCOUNT_TYPE_SSL)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslUpdateCertificateRecordRetry,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ssl_update_certificate_record_retry_operation.update_certificate_record_retry", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ssl_update_certificate_record_retry_operation.update_certificate_record_retry", "deploy_record_id", "1689"),
				),
			},
			{
				ResourceName:      "tencentcloud_ssl_update_certificate_record_retry_operation.update_certificate_record_retry",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSslUpdateCertificateRecordRetry = `

resource "tencentcloud_ssl_update_certificate_record_retry_operation" "update_certificate_record_retry" {
  deploy_record_id = "1689"
}

`
