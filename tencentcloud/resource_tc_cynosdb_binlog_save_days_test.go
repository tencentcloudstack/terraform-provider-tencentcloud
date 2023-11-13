package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

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
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_binlog_save_days.binlog_save_days", "id")),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_binlog_save_days.binlog_save_days",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCynosdbBinlogSaveDays = `

resource "tencentcloud_cynosdb_binlog_save_days" "binlog_save_days" {
  cluster_id = "cynosdbmysql-123"
  binlog_save_days = 7
}

`
