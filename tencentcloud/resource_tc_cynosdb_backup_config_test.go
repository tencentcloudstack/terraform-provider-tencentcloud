package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbBackupConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbBackupConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_backup_config.backup_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_backup_config.backup_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCynosdbBackupConfig = `

resource "tencentcloud_cynosdb_backup_config" "backup_config" {
  cluster_id = &lt;nil&gt;
  backup_time_beg = 3600
  backup_time_end = &lt;nil&gt;
  reserve_duration = &lt;nil&gt;
  backup_freq = &lt;nil&gt;
  backup_type = &lt;nil&gt;
}

`
