package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudClickhouseBackupStrategyResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClickhouseBackupStrategy,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_clickhouse_backup_strategy.backup_strategy", "id")),
			},
			{
				ResourceName:      "tencentcloud_clickhouse_backup_strategy.backup_strategy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccClickhouseBackupStrategy = DefaultClickhouseVariables + `
resource "tencentcloud_clickhouse_backup" "backup" {
	instance_id = var.instance_id
	cos_bucket_name = "keep-export-image-1308726196"
}

resource "tencentcloud_clickhouse_backup_strategy" "backup_strategy" {
	instance_id = var.instance_id
	data_backup_strategy {
	  week_days = "3"
	  retain_days = 2
	  execute_hour = 1
	  back_up_tables {
		database = "iac"
		table = "my_table"
		total_bytes = 0
		v_cluster = "default_cluster"
		ips = "10.0.0.35"
	  }
	}
	meta_backup_strategy {
		week_days = "1"
		retain_days = 2
		execute_hour = 3
	}
}
`
