package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const TestObjectPgModifySwitchTimePeriod = "tencentcloud_postgresql_modify_switch_time_period_operation.modify_switch_time_period_operation"

// For now, we only supports updating the minor kernel version immediately, which mean is switch_tag=0
func TestAccTencentCloudNeedFixPostgresqlModifySwitchTimePeriodOperationResource_basic(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccStepSetRegion(t, "ap-guangzhou")
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlModifySwitchTimePeriodOperation,
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-guangzhou")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(TestObjectPgModifySwitchTimePeriod, "id"),
					resource.TestCheckResourceAttrSet(TestObjectPgModifySwitchTimePeriod, "db_instance_id"),
					resource.TestCheckResourceAttr(TestObjectPgModifySwitchTimePeriod, "switch_tag", "0"),
				),
			},
		},
	})
}

const testAccPostgresqlModifySwitchTimePeriodOperation = OperationPresetPGSQL + `

resource "tencentcloud_postgresql_modify_switch_time_period_operation" "modify_switch_time_period_operation" {
  db_instance_id = local.pgsql_id
  switch_tag = 0
}

`
