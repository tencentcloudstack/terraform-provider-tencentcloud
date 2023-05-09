package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dbbrain_tdsql_audit_log.my_log", "id")),
			},
			{
				ResourceName:      "tencentcloud_dbbrain_tdsql_audit_log.my_log",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDbbrainTdsqlAuditLog = `

resource "tencentcloud_dbbrain_tdsql_audit_log" "my_log" {
  product = ""
  node_request_type = ""
  instance_id = ""
  start_time = ""
  end_time = ""
  filter {
		host = 
		db_name = 
		user = 
		sent_rows = 
		affect_rows = 
		exec_time = 

  }
}

`
