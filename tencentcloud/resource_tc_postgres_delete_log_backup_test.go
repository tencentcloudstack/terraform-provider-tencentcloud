package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPostgresDeleteLogBackupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresDeleteLogBackup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_postgres_delete_log_backup.delete_log_backup", "id")),
			},
			{
				ResourceName:      "tencentcloud_postgres_delete_log_backup.delete_log_backup",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPostgresDeleteLogBackup = `

resource "tencentcloud_postgres_delete_log_backup" "delete_log_backup" {
  d_b_instance_id = ""
  log_backup_id = ""
  tags = {
    "createdBy" = "terraform"
  }
}

`
