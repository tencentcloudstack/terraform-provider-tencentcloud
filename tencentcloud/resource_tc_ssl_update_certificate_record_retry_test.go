package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSslUpdateCertificateRecordRetryResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslUpdateCertificateRecordRetry,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ssl_update_certificate_record_retry.update_certificate_record_retry", "id")),
			},
			{
				ResourceName:      "tencentcloud_ssl_update_certificate_record_retry.update_certificate_record_retry",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSslUpdateCertificateRecordRetry = `

resource "tencentcloud_ssl_update_certificate_record_retry" "update_certificate_record_retry" {
  deploy_record_id = 
  deploy_record_detail_id = 
}

`
