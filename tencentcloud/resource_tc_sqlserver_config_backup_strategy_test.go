package tencentcloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudSqlserverConfigBackupStrategyResource_basic(t *testing.T) {
	t.Parallel()
	loc, _ := time.LoadLocation("Asia/Chongqing")
	startTime := time.Now().In(loc).Format("2006-01-02")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverConfigBackupStrategy_daily,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_config_backup_strategy.config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "backup_type", "daily"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "backup_time", "0"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "backup_day", "1"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "backup_model", "master_no_pkg"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "backup_save_days", "7"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "regular_backup_enable", "disable"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "regular_backup_save_days", "90"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "regular_backup_strategy", "months"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "regular_backup_counts", "1"),
				),
			},
			{
				Config: testAccSqlserverConfigBackupStrategy_weekly,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_config_backup_strategy.config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "backup_type", "weekly"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "backup_time", "0"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "backup_day", "1"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "backup_model", "master_no_pkg"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "backup_cycle.#", "3"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "backup_save_days", "7"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "regular_backup_enable", "disable"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "regular_backup_save_days", "90"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "regular_backup_strategy", "months"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "regular_backup_counts", "1"),
				),
			},
			{
				Config: fmt.Sprintf(testAccSqlserverConfigBackupStrategy_regular_months, startTime),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_config_backup_strategy.config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "backup_time", "0"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "backup_day", "1"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "backup_model", "master_no_pkg"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "backup_cycle.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "backup_save_days", "7"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "regular_backup_enable", "enable"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "regular_backup_save_days", "120"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "regular_backup_strategy", "months"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "regular_backup_counts", "1"),
				),
			},
			{
				Config: fmt.Sprintf(testAccSqlserverConfigBackupStrategy_regular_yearly, startTime),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_config_backup_strategy.config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "backup_time", "0"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "backup_day", "1"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "backup_model", "master_no_pkg"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "backup_cycle.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "backup_save_days", "7"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "regular_backup_enable", "enable"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "regular_backup_save_days", "1000"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "regular_backup_strategy", "years"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_config_backup_strategy.config", "regular_backup_counts", "2"),
				),
			},
			{
				ResourceName: "tencentcloud_sqlserver_config_backup_strategy.config",
				ImportState:  true,
			},
		},
	})
}

const testAccSqlserverConfigBackupStrategy_daily = CommonPresetSQLServer + `

resource "tencentcloud_sqlserver_config_backup_strategy" "config" {
  instance_id = local.sqlserver_id
  backup_type = "daily"
  backup_time = 0
  backup_day = 1
  backup_model = "master_no_pkg"
  backup_cycle = [1]
  backup_save_days = 7
  regular_backup_enable = "disable"
  regular_backup_save_days = 90
  regular_backup_strategy = "months"
  regular_backup_counts = 1
}

`

const testAccSqlserverConfigBackupStrategy_weekly = CommonPresetSQLServer + `

resource "tencentcloud_sqlserver_config_backup_strategy" "config" {
  instance_id = local.sqlserver_id
  backup_type = "weekly"
  backup_time = 0
  backup_day = 1
  backup_model = "master_no_pkg"
  backup_cycle = [1,3,5]
  backup_save_days = 7
  regular_backup_enable = "disable"
  regular_backup_save_days = 90
  regular_backup_strategy = "months"
  regular_backup_counts = 1
}

`

const testAccSqlserverConfigBackupStrategy_regular_months = CommonPresetSQLServer + `

resource "tencentcloud_sqlserver_config_backup_strategy" "config" {
  instance_id = local.sqlserver_id
  backup_type = "weekly"
  backup_time = 0
  backup_day = 1
  backup_model = "master_no_pkg"
  backup_cycle = [1,3]
  backup_save_days = 7
  regular_backup_enable = "enable"
  regular_backup_save_days = 120
  regular_backup_strategy = "months"
  regular_backup_counts = 1
  regular_backup_start_time = "%s"
}

`

const testAccSqlserverConfigBackupStrategy_regular_yearly = CommonPresetSQLServer + `

resource "tencentcloud_sqlserver_config_backup_strategy" "config" {
  instance_id = local.sqlserver_id
  backup_type = "weekly"
  backup_time = 0
  backup_day = 1
  backup_model = "master_no_pkg"
  backup_cycle = [1,3]
  backup_save_days = 7
  regular_backup_enable = "enable"
  regular_backup_save_days = 1000
  regular_backup_strategy = "years"
  regular_backup_counts = 2
  regular_backup_start_time = "%s"
}

`
