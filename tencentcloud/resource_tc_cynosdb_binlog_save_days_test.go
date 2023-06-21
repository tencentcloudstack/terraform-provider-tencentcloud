package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbBinlogSaveDaysResource_basic -v
func TestAccTencentCloudCynosdbBinlogSaveDaysResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbBinlogSaveDays,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_binlog_save_days.binlog_save_days", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_binlog_save_days.binlog_save_days", "binlog_save_days", "8"),
				),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_binlog_save_days.binlog_save_days",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCynosdbBinlogSaveDaysUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_binlog_save_days.binlog_save_days", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_binlog_save_days.binlog_save_days", "binlog_save_days", "7"),
				),
			},
		},
	})
}

const testAccCynosdbBinlogSaveDays = CommonCynosdb + `

resource "tencentcloud_cynosdb_binlog_save_days" "binlog_save_days" {
  cluster_id = var.cynosdb_cluster_id
  binlog_save_days = 8
}

`

const testAccCynosdbBinlogSaveDaysUp = CommonCynosdb + `

resource "tencentcloud_cynosdb_binlog_save_days" "binlog_save_days" {
  cluster_id = var.cynosdb_cluster_id
  binlog_save_days = 7
}

`
