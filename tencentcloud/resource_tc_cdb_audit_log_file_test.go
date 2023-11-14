package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbAuditLogFileResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbAuditLogFile,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cdb_audit_log_file.audit_log_file", "id")),
			},
			{
				ResourceName:      "tencentcloud_cdb_audit_log_file.audit_log_file",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCdbAuditLogFile = `

resource "tencentcloud_cdb_audit_log_file" "audit_log_file" {
  instance_id = "cdb-c1nl9rpv"
  start_time = "2022-07-12 10:29:20"
  end_time = "2022-08-12 10:29:20"
  order = "ASC"
  order_by = ""
  filter {
		host = &lt;nil&gt;
		user = &lt;nil&gt;
		d_b_name = &lt;nil&gt;
		table_name = &lt;nil&gt;
		policy_name = &lt;nil&gt;
		sql = &lt;nil&gt;
		sql_type = &lt;nil&gt;
		exec_time = &lt;nil&gt;
		affect_rows = &lt;nil&gt;
		sql_types = &lt;nil&gt;
		sqls = &lt;nil&gt;

  }
    }

`
