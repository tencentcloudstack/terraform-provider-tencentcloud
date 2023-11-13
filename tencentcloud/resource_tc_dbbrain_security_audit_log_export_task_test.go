package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDbbrainSecurityAuditLogExportTaskResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDbbrainSecurityAuditLogExportTask,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dbbrain_security_audit_log_export_task.security_audit_log_export_task", "id")),
			},
			{
				ResourceName:      "tencentcloud_dbbrain_security_audit_log_export_task.security_audit_log_export_task",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDbbrainSecurityAuditLogExportTask = `

resource "tencentcloud_dbbrain_security_audit_log_export_task" "security_audit_log_export_task" {
  sec_audit_group_id = &lt;nil&gt;
  start_time = &lt;nil&gt;
  end_time = &lt;nil&gt;
  product = &lt;nil&gt;
  danger_levels = &lt;nil&gt;
}

`
