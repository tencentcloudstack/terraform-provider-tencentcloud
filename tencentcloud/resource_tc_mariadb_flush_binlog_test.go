package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMariadbFlushBinlogResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbFlushBinlog,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mariadb_flush_binlog.flush_binlog", "id")),
			},
			{
				ResourceName:      "tencentcloud_mariadb_flush_binlog.flush_binlog",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMariadbFlushBinlog = `

resource "tencentcloud_mariadb_flush_binlog" "flush_binlog" {
  instance_id = ""
}

`
