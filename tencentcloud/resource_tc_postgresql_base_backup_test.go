package tencentcloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccPostgresqlBaseBackupObject = "tencentcloud_postgresql_base_backup.base_backup"

func TestAccTencentCloudPostgresqlBaseBackupResource_basic(t *testing.T) {
	// t.Parallel()
	newExpireTime := time.Now().AddDate(0, 0, +1).Format("2006-01-02 15:04:05")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccStepSetRegion(t, "ap-chengdu")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlBaseBackup,
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-chengdu")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupObject, "id"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupObject, "db_instance_id"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupObject, "base_backup_id"),
				),
			},
			{
				Config: fmt.Sprintf(testAccPostgresqlBaseBackup_update, newExpireTime),
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-chengdu")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
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

const testAccPostgresqlInstance_common string = `
data "tencentcloud_security_groups" "cd" {
  name = "keep-tf-test"
}

data "tencentcloud_postgresql_instances" "cd" {
  name = "keep_postsql_instance"
}

data "tencentcloud_vpc_subnets" "cd1" {
  name = "keep-tf-test-subnet"
}

locals {
  sg_id      = data.tencentcloud_security_groups.cd.security_groups.0.security_group_id
  vpc_id     = data.tencentcloud_vpc_subnets.cd1.instance_list.0.vpc_id
  subnet_id  = data.tencentcloud_vpc_subnets.cd1.instance_list.0.subnet_id
  pgsql_id   = data.tencentcloud_postgresql_instances.cd.instance_list.0.id
  default_az = "ap-chengdu-1"
}


`

const testAccPostgresqlBaseBackup = testAccPostgresqlInstance_common + `

resource "tencentcloud_postgresql_base_backup" "base_backup" {
  db_instance_id = local.pgsql_id
  tags = {
    "createdBy" = "terraform"
  }
}

`

const testAccPostgresqlBaseBackup_update = testAccPostgresqlInstance_common + `

resource "tencentcloud_postgresql_base_backup" "base_backup" {
  db_instance_id = local.pgsql_id
  new_expire_time = "%s"
  tags = {
    "createdBy" = "terraform"
  }
}

`
