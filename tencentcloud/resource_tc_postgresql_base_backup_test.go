package tencentcloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccPostgresqlBaseBackupObject = "tencentcloud_postgresql_base_backup.base_backup"

func TestAccTencentCloudPostgresqlBaseBackupResource_basic(t *testing.T) {
	t.Parallel()
	newExpireTime := time.Now().AddDate(0, 0, +1).Format("2006-01-02 15:04:05")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlBaseBackup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupObject, "id"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupObject, "db_instance_id"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupObject, "base_backup_id"),
				),
			},
			{
				Config: fmt.Sprintf(testAccPostgresqlBaseBackup_update, newExpireTime),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupObject, "id"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupObject, "db_instance_id"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupObject, "base_backup_id"),
					resource.TestCheckResourceAttr(testAccPostgresqlBaseBackupObject, "new_expire_time", newExpireTime),
				),
			},
		},
	})
}

const testAccPostgresqlBaseBackup = CommonPresetPGSQL + `

resource "tencentcloud_postgresql_base_backup" "base_backup" {
  db_instance_id = local.pgsql_id
  tags = {
    "createdBy" = "terraform"
  }
}

`

const testAccPostgresqlBaseBackup_update = CommonPresetPGSQL + `

resource "tencentcloud_postgresql_base_backup" "base_backup" {
  db_instance_id = local.pgsql_id
  new_expire_time = "%s"
  tags = {
    "createdBy" = "terraform"
  }
}

`
