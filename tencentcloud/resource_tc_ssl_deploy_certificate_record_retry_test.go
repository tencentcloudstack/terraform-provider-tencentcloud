package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
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
			{
				ResourceName:      "tencentcloud_ssl_deploy_certificate_record_retry.deploy_certificate_record_retry",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSslDeployCertificateRecordRetry = `

resource "tencentcloud_ssl_deploy_certificate_record_retry" "deploy_certificate_record_retry" {
  deploy_record_id = 
  deploy_record_detail_id = 
}

`
