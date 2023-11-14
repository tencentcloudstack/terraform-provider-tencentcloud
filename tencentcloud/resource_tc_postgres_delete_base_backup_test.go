package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPostgresDeleteBaseBackupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresDeleteBaseBackup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_postgres_delete_base_backup.delete_base_backup", "id")),
			},
			{
				ResourceName:      "tencentcloud_postgres_delete_base_backup.delete_base_backup",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPostgresDeleteBaseBackup = `

resource "tencentcloud_postgres_delete_base_backup" "delete_base_backup" {
  d_b_instance_id = ""
  base_backup_id = ""
  tags = {
    "createdBy" = "terraform"
  }
}

`
