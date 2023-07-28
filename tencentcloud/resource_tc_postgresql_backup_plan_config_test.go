package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccPostgresqlBackupPlanConfigObject = "tencentcloud_postgresql_backup_plan_config.backup_plan_config"

func TestAccTencentCloudPostgresqlBackupPlanConfigResource_basic(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccStepSetRegion(t, "ap-guangzhou")
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresBackupPlanConfig,
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-guangzhou")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAccPostgresqlBackupPlanConfigObject, "id"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBackupPlanConfigObject, "db_instance_id"),
					resource.TestCheckResourceAttr(testAccPostgresqlBackupPlanConfigObject, "min_backup_start_time", "01:00:00"),
					resource.TestCheckResourceAttr(testAccPostgresqlBackupPlanConfigObject, "max_backup_start_time", "02:00:00"),
					resource.TestCheckResourceAttr(testAccPostgresqlBackupPlanConfigObject, "base_backup_retention_period", "7"),
					resource.TestCheckResourceAttr(testAccPostgresqlBackupPlanConfigObject, "backup_period.#", "3"),
					resource.TestCheckTypeSetElemAttr(testAccPostgresqlBackupPlanConfigObject, "backup_period.*", "monday"),
					resource.TestCheckTypeSetElemAttr(testAccPostgresqlBackupPlanConfigObject, "backup_period.*", "wednesday"),
					resource.TestCheckTypeSetElemAttr(testAccPostgresqlBackupPlanConfigObject, "backup_period.*", "friday"),
				),
			},
			{
				ResourceName:      testAccPostgresqlBackupPlanConfigObject,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPostgresBackupPlanConfig = CommonPresetPGSQL + `

resource "tencentcloud_postgresql_backup_plan_config" "backup_plan_config" {
  db_instance_id = local.pgsql_id
  min_backup_start_time = "01:00:00"
  max_backup_start_time = "02:00:00"
  base_backup_retention_period = 7
  backup_period = ["monday","wednesday","friday"]
}

`
