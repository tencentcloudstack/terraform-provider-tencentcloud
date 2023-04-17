package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_config_backup_strategy.config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "backup_type", "weekly"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "backup_time", "0"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "backup_day", "1"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "backup_model", "master_no_pkg"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "backup_cycle.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "backup_save_days", "7"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "regular_backup_enable", "disable"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "regular_backup_save_days", "365"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "regular_backup_strategy", "months"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "regular_backup_counts", "1"),
				),
			},
			{
				ResourceName: "tencentcloud_sqlserver_config_backup_strategy.config",
				ImportState:  true,
			},
		},
	})
}

const testAccSqlserverConfigBackupStrategy = CommonPresetSQLServer + `

resource "tencentcloud_sqlserver_config_backup_strategy" "config" {
  instance_id = local.sqlserver_id
  backup_type = "weekly"
  backup_time = 0
  backup_day = 1
  backup_model = "master_no_pkg"
  backup_cycle = [1,4]
  backup_save_days = 7
  regular_backup_enable = "disable"
  regular_backup_save_days = 365
  regular_backup_strategy = "months"
  regular_backup_counts = 1
}

`
