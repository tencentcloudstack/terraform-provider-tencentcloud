package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbExportInstanceErrorLogsResource_basic -v
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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_export_instance_error_logs.export_instance_error_logs", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_export_instance_error_logs.export_instance_error_logs", "error_log_item_export.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_export_instance_error_logs.export_instance_error_logs", "error_log_item_export.0.timestamp"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_export_instance_error_logs.export_instance_error_logs", "error_log_item_export.0.level"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_export_instance_error_logs.export_instance_error_logs", "error_log_item_export.0.content"),
				),
			},
		},
	})
}

const testAccCynosdbExportInstanceErrorLogs = CommonCynosdb + `

resource "tencentcloud_cynosdb_export_instance_error_logs" "export_instance_error_logs" {
  instance_id = var.cynosdb_cluster_instance_id
}

`
