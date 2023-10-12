package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslDeployCertificateRecordRollbackResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCommon(t, ACCOUNT_TYPE_SSL)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslDeployCertificateRecordRollback,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ssl_deploy_certificate_record_rollback_operation.deploy_certificate_record_rollback", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ssl_deploy_certificate_record_rollback_operation.deploy_certificate_record_rollback", "deploy_record_id", "36058"),
				),
			},
		},
	})
}

const testAccSslDeployCertificateRecordRollback = `

resource "tencentcloud_ssl_deploy_certificate_record_rollback_operation" "deploy_certificate_record_rollback" {
  deploy_record_id = 36058
}

`
