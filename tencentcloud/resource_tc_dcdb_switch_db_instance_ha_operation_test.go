package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDcdbSwitchDbInstanceHaOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbSwitchDbInstanceHaOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_switch_db_instance_ha_operation.switch_operation", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_switch_db_instance_ha_operation.switch_operation", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_switch_db_instance_ha_operation.switch_operation", "zone"),
				),
			},
			{
				Config: testAccDcdbSwitchDbInstanceHaOperation_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_switch_db_instance_ha_operation.switch_operation", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_switch_db_instance_ha_operation.switch_operation", "instance_id"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_switch_db_instance_ha_operation.switch_operation", "zone", "ap-guangzhou-3"),
				),
			},
		},
	})
}

const testAccDcdbSwitchDbInstanceHaOperation = CommonPresetDcdb + `

resource "tencentcloud_dcdb_switch_db_instance_ha_operation" "switch_operation" {
  instance_id = local.dcdb_id
  zone = "ap-guangzhou-4" //3 to 4
}

`

const testAccDcdbSwitchDbInstanceHaOperation_update = CommonPresetDcdb + `

resource "tencentcloud_dcdb_switch_db_instance_ha_operation" "switch_operation" {
  instance_id = local.dcdb_id
  zone = "ap-guangzhou-3" //4 to 3
}

`
