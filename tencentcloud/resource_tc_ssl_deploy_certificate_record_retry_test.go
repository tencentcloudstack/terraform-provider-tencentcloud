package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslDeployCertificateRecordRetryResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslDeployCertificateRecordRetry,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ssl_deploy_certificate_record_retry.deploy_certificate_record_retry", "id")),
			},
		},
	})
}

const testAccSslDeployCertificateRecordRetry = `

resource "tencentcloud_ssl_deploy_certificate_record_retry" "deploy_certificate_record_retry" {
  deploy_record_id = 35521
}

`
