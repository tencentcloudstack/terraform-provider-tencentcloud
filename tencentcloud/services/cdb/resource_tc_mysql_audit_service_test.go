package cdb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudMysqlAuditServiceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlAuditService,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_audit_service.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_audit_service.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_audit_service.example", "log_expire_day"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_audit_service.example", "high_log_expire_day"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_audit_service.example", "rule_template_ids"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_audit_service.example", "audit_all"),
				),
			},
			{
				Config: testAccMysqlAuditServiceUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_audit_service.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_audit_service.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_audit_service.example", "log_expire_day"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_audit_service.example", "high_log_expire_day"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_audit_service.example", "audit_all"),
				),
			},
			{
				ResourceName:      "tencentcloud_mysql_audit_service.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMysqlAuditService = `
resource "tencentcloud_mysql_audit_service" "example" {
  instance_id         = "cdb-3kwa3gfj"
  log_expire_day      = 90
  high_log_expire_day = 30
  rule_template_ids   = [
    "cdb-art-3a9ww0oj"
  ]
  audit_all           = false
}
`

const testAccMysqlAuditServiceUpdate = `
resource "tencentcloud_mysql_audit_service" "example" {
  instance_id         = "cdb-3kwa3gfj"
  log_expire_day      = 30
  high_log_expire_day = 7
  audit_all           = true
}
`
