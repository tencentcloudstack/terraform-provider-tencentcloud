package cynosdb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCynosdbBackupConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbBackupConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_backup_config.foo", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_backup_config.foo", "backup_time_beg", "7200"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_backup_config.foo", "backup_time_end", "21600"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_backup_config.foo", "reserve_duration", "604800"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_backup_config.foo", "logic_backup_config.0.logic_backup_enable", "ON"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_backup_config.foo", "logic_backup_config.0.logic_backup_time_beg", "7200"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_backup_config.foo", "logic_backup_config.0.logic_backup_time_end", "21600"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_backup_config.foo", "logic_backup_config.0.logic_cross_regions.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_backup_config.foo", "logic_backup_config.0.logic_cross_regions.0", "ap-shanghai"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_backup_config.foo", "logic_backup_config.0.logic_cross_regions_enable", "ON"),
				),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_backup_config.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCynosdbBackupConfigUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_backup_config.foo", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_backup_config.foo", "backup_time_beg", "7200"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_backup_config.foo", "backup_time_end", "21600"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_backup_config.foo", "reserve_duration", "604800"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_backup_config.foo", "logic_backup_config.0.logic_backup_enable", "OFF"),
				),
			},
		},
	})
}

const testAccCynosdbBackupConfig = testAccCynosdbCluster + `

resource "tencentcloud_cynosdb_backup_config" "foo" {
    backup_time_beg  = 7200
    backup_time_end  = 21600
    cluster_id       = tencentcloud_cynosdb_cluster.foo.id
    reserve_duration = 604800

    logic_backup_config {
        logic_backup_enable        = "ON"
        logic_backup_time_beg      = 7200
        logic_backup_time_end      = 21600
        logic_cross_regions        = ["ap-shanghai"]
        logic_cross_regions_enable = "ON"
        logic_reserve_duration     = 259200
    }
}
`

const testAccCynosdbBackupConfigUp = testAccCynosdbCluster + `

resource "tencentcloud_cynosdb_backup_config" "foo" {
    backup_time_beg  = 7200
    backup_time_end  = 21600
    cluster_id       = tencentcloud_cynosdb_cluster.foo.id
    reserve_duration = 604800

    logic_backup_config {
        logic_backup_enable        = "OFF"
    }
}
`
