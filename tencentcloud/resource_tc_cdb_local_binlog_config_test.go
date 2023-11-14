package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbLocalBinlogConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbLocalBinlogConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cdb_local_binlog_config.local_binlog_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_cdb_local_binlog_config.local_binlog_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCdbLocalBinlogConfig = `

resource "tencentcloud_cdb_local_binlog_config" "local_binlog_config" {
  instance_id = ""
  save_hours = 
  max_usage = 
}

`
