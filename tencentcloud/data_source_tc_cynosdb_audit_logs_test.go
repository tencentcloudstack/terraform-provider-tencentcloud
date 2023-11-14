package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbAuditLogsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbAuditLogsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_audit_logs.audit_logs")),
			},
		},
	})
}

const testAccCynosdbAuditLogsDataSource = `

data "tencentcloud_cynosdb_audit_logs" "audit_logs" {
  instance_id = "cynosdbmysql-ins-xx"
  start_time = "2017-07-12 10:29:20"
  end_time = "2017-07-12 10:29:20"
  order = "ASC"
  order_by = "timestap"
  filter {
		host = 
		user = 
		d_b_name = 
		table_name = 
		policy_name = 
		sql = ""
		sql_type = ""
		exec_time = 
		affect_rows = 
		sql_types = 
		sqls = 
		sent_rows = 
		thread_id = 

  }
  }

`
