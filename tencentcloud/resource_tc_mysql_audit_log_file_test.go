package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudMysqlAuditLogFileResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlAuditLogFile,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mysql_audit_log_file.audit_log_file", "id")),
			},
		},
	})
}

const testAccMysqlAuditLogFile = `

resource "tencentcloud_mysql_audit_log_file" "audit_log_file" {
  instance_id = "cdb-fitq5t9h"
  start_time  = "2023-03-28 20:14:00"
  end_time    = "2023-03-29 20:14:00"
  order       = "ASC"
  order_by    = "timestamp"
  filter {
    host = ["30.50.207.46"]
    user = ["keep_dbbrain"]
  }
}

`
