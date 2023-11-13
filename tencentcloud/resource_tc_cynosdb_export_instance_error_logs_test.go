package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbExportInstanceErrorLogsResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbExportInstanceErrorLogs,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_export_instance_error_logs.export_instance_error_logs", "id")),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_export_instance_error_logs.export_instance_error_logs",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCynosdbExportInstanceErrorLogs = `

resource "tencentcloud_cynosdb_export_instance_error_logs" "export_instance_error_logs" {
  instance_id = "cynosdbmysql-ins-123"
  start_time = "2022-01-01 12:00:00"
  end_time = "2022-01-01 14:00:00"
  log_levels = 
  key_words = 
  file_type = "csv"
  order_by = "Timestamp"
  order_by_type = "ASC"
}

`
