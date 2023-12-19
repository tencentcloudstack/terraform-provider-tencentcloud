package cdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlBackupEncryptionStatusResource_basic -v
func TestAccTencentCloudMysqlBackupEncryptionStatusResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlBackupEncryptionStatus,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_backup_encryption_status.backup_encryption_status", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_backup_encryption_status.backup_encryption_status", "encryption_status", "on"),
				),
			},
			{
				ResourceName:      "tencentcloud_mysql_backup_encryption_status.backup_encryption_status",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMysqlBackupEncryptionStatusUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_backup_encryption_status.backup_encryption_status", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_backup_encryption_status.backup_encryption_status", "encryption_status", "off"),
				),
			},
		},
	})
}

const testAccMysqlBackupEncryptionStatusVar = `
variable "instance_id" {
  default = "` + tcacctest.DefaultDbBrainInstanceId + `"
}
`

const testAccMysqlBackupEncryptionStatus = testAccMysqlBackupEncryptionStatusVar + `

resource "tencentcloud_mysql_backup_encryption_status" "backup_encryption_status" {
	instance_id = var.instance_id
	encryption_status = "on"
}

`

const testAccMysqlBackupEncryptionStatusUp = testAccMysqlBackupEncryptionStatusVar + `

resource "tencentcloud_mysql_backup_encryption_status" "backup_encryption_status" {
	instance_id = var.instance_id
	encryption_status = "off"
}

`
