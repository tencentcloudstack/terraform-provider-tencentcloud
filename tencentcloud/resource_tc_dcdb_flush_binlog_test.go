package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDcdbFlushBinlogResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbFlushBinlog,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dcdb_flush_binlog.flush_binlog", "id")),
			},
			{
				ResourceName:      "tencentcloud_dcdb_flush_binlog.flush_binlog",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDcdbFlushBinlog = `

resource "tencentcloud_dcdb_flush_binlog" "flush_binlog" {
  instance_id = ""
}

`
