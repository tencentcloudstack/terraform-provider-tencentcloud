package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDcdbFlushBinlogOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbFlushBinlogOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_flush_binlog_operation.flush_operation", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_flush_binlog_operation.flush_operation", "instance_id"),
				),
			},
		},
	})
}

const testAccDcdbFlushBinlogOperation = CommonPresetDcdb + `

resource "tencentcloud_dcdb_flush_binlog_operation" "flush_operation" {
  instance_id = local.dcdb_id
}

`
