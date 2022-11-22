package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudDbbrainSecurityAuditLogExportTask_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDbbrainSecurityAuditLogExportTask,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dbbrain_security_audit_log_export_task.security_audit_log_export_task", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_dbbrain_security_audit_log_export_task.securityAuditLogExportTask",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDbbrainSecurityAuditLogExportTask = `

resource "tencentcloud_dbbrain_security_audit_log_export_task" "security_audit_log_export_task" {
  sec_audit_group_id = ""
  start_time = ""
  end_time = ""
  product = ""
  danger_levels = ""
}

`
