package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverConfigBackupStrategyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverConfigBackupStrategy,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_config_backup_strategy.config_backup_strategy", "id")),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_config_backup_strategy.config_backup_strategy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverConfigBackupStrategy = `

resource "tencentcloud_sqlserver_config_backup_strategy" "config_backup_strategy" {
  instance_id = "mssql-i1z41iwd"
  backup_type = "weekly"
  backup_time = 0
  backup_day = 1
  backup_model = "master_pkg"
  backup_cycle = 
  backup_save_days = 10
  regular_backup_enable = "enabled"
  regular_backup_save_days = 365
  regular_backup_strategy = "monthly"
  regular_backup_counts = 3
  regular_backup_start_time = "2023-04-10"
}

`
