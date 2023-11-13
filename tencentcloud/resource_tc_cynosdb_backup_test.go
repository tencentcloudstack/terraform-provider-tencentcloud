package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbBackupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbBackup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_backup.backup", "id")),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_backup.backup",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCynosdbBackup = `

resource "tencentcloud_cynosdb_backup" "backup" {
  cluster_id = &lt;nil&gt;
  backup_type = &lt;nil&gt;
  backup_databases = &lt;nil&gt;
  backup_tables {
		database = ""
		tables = &lt;nil&gt;

  }
  backup_name = &lt;nil&gt;
}

`
