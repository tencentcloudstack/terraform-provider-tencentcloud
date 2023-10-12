package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslDeployCertificateRecordRetryResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCommon(t, ACCOUNT_TYPE_SSL)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslDeployCertificateRecordRetry,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ssl_deploy_certificate_record_retry_operation.deploy_certificate_record_retry", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ssl_deploy_certificate_record_retry_operation.deploy_certificate_record_retry", "deploy_record_id", "36062"),
				),
			},
		},
	})
}

const testAccSslDeployCertificateRecordRetry = `

resource "tencentcloud_ssl_deploy_certificate_record_retry_operation" "deploy_certificate_record_retry" {
  deploy_record_id = 36062
}

`
