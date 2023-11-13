package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbBackupEncryptionStatusResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbBackupEncryptionStatus,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cdb_backup_encryption_status.backup_encryption_status", "id")),
			},
			{
				ResourceName:      "tencentcloud_cdb_backup_encryption_status.backup_encryption_status",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCdbBackupEncryptionStatus = `

resource "tencentcloud_cdb_backup_encryption_status" "backup_encryption_status" {
  instance_id = "cdb-c1nl9rpv"
  encryption_status = "on"
}

`
