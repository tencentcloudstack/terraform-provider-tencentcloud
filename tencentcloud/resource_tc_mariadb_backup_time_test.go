package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMariadbBackupTimeResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbBackupTime,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mariadb_backup_time.backup_time", "id")),
			},
			{
				ResourceName:      "tencentcloud_mariadb_backup_time.backup_time",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMariadbBackupTime = `

resource "tencentcloud_mariadb_backup_time" "backup_time" {
  instance_id = ""
  start_backup_time = ""
  end_backup_time = ""
}

`
