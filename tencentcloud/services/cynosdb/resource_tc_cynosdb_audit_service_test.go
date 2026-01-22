package cynosdb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCynosdbAuditServiceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbAuditService,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_audit_service.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_audit_service.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_audit_service.example", "log_expire_day"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_audit_service.example", "high_log_expire_day"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_audit_service.example", "audit_all"),
				),
			},
			{
				Config: testAccCynosdbAuditServiceUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_audit_service.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_audit_service.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_audit_service.example", "log_expire_day"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_audit_service.example", "high_log_expire_day"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_audit_service.example", "rule_template_ids"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_audit_service.example", "audit_all"),
				),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_audit_service.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCynosdbAuditService = `
resource "tencentcloud_cynosdb_audit_service" "example" {
  instance_id         = "cynosdbmysql-31zv4ii1"
  log_expire_day      = 30
  high_log_expire_day = 7
  audit_all           = true
}
`

const testAccCynosdbAuditServiceUpdate = `
resource "tencentcloud_cynosdb_audit_service" "example" {
  instance_id         = "cynosdbmysql-31zv4ii1"
  log_expire_day      = 30
  high_log_expire_day = 7
  rule_template_ids   = ["cynosdb-art-riwq2vx0"]
  audit_all           = false
}
`
