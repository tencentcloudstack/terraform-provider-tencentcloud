package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccTencentCloudDbbrainTdsqlAuditLogResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDbbrainTdsqlAuditLog,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dbbrain_tdsql_audit_log.tdsql_audit_log", "id")),
			},
			{
				ResourceName:      "tencentcloud_dbbrain_tdsql_audit_log.tdsql_audit_log",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDbbrainTdsqlAuditLog = `

resource "tencentcloud_dbbrain_tdsql_audit_log" "tdsql_audit_log" {
  product = ""
  node_request_type = ""
  instance_id = ""
  start_time = ""
  end_time = ""
  filter {
		host = 
		d_b_name = 
		user = 
		sent_rows = 
		affect_rows = 
		exec_time = 

  }
}

`
