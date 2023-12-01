package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbAuditLogsDataSource_basic -v
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_audit_logs.audit_logs"),
				),
			},
		},
	})
}

const testAccCynosdbAuditLogsDataSource = `
data "tencentcloud_cynosdb_audit_logs" "audit_logs" {
  instance_id = "cynosdbmysql-ins-afqx1hy0"
  start_time  = "2023-06-18 10:00:00"
  end_time    = "2023-06-18 10:00:02"
  order       = "DESC"
  order_by    = "timestamp"
  filter {
    host        = ["30.50.207.176"]
    user        = ["keep_dts"]
    policy_name = ["default_audit"]
    sql_type    = "SELECT"
  }
}
`
