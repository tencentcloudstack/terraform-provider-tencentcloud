package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudMysqlLocalBinlogConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlLocalBinlogConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mysql_local_binlog_config.local_binlog_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_mysql_local_binlog_config.local_binlog_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMysqlLocalBinlogConfig = `

resource "tencentcloud_mysql_local_binlog_config" "local_binlog_config" {
  instance_id = "cdb-fitq5t9h"
  save_hours = 140
  max_usage = 50
}

`
