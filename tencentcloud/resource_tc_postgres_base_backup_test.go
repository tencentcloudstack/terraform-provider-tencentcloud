package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPostgresBaseBackupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresBaseBackup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_postgres_base_backup.base_backup", "id")),
			},
			{
				ResourceName:      "tencentcloud_postgres_base_backup.base_backup",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPostgresBaseBackup = `

resource "tencentcloud_postgres_base_backup" "base_backup" {
  d_b_instance_id = ""
  base_backup_id = ""
  new_expire_time = ""
}

`
