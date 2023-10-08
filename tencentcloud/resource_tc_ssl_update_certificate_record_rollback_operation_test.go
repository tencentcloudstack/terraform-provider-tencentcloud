package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslUpdateCertificateRecordRollbackResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCommon(t, ACCOUNT_TYPE_SSL)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslUpdateCertificateRecordRollback,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ssl_update_certificate_record_rollback_operation.update_certificate_record_rollback", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ssl_update_certificate_record_rollback_operation.update_certificate_record_rollback", "deploy_record_id", "1693"),
				),
			},
		},
	})
}

const testAccSslUpdateCertificateRecordRollback = `

resource "tencentcloud_ssl_update_certificate_record_rollback_operation" "update_certificate_record_rollback" {
  deploy_record_id = "1693"
}

`
