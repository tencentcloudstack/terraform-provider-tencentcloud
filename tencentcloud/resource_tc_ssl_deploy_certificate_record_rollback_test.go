package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSslDeployCertificateRecordRollbackResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslDeployCertificateRecordRollback,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ssl_deploy_certificate_record_rollback.deploy_certificate_record_rollback", "id")),
			},
			{
				ResourceName:      "tencentcloud_ssl_deploy_certificate_record_rollback.deploy_certificate_record_rollback",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSslDeployCertificateRecordRollback = `

resource "tencentcloud_ssl_deploy_certificate_record_rollback" "deploy_certificate_record_rollback" {
  deploy_record_id = 
}

`
