package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDbbrainDbDiagReportTaskResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDbbrainDbDiagReportTask,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dbbrain_db_diag_report_task.db_diag_report_task", "id")),
			},
			{
				ResourceName:      "tencentcloud_dbbrain_db_diag_report_task.db_diag_report_task",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDbbrainDbDiagReportTask = `

resource "tencentcloud_dbbrain_db_diag_report_task" "db_diag_report_task" {
  instance_id = ""
  start_time = ""
  end_time = ""
  send_mail_flag = 
  contact_person = 
  contact_group = 
  product = ""
}

`
