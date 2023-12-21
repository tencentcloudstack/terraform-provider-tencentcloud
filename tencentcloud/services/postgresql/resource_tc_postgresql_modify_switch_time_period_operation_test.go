package postgresql_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const TestObjectPgModifySwitchTimePeriod = "tencentcloud_postgresql_modify_switch_time_period_operation.modify_switch_time_period_operation"

// For now, we only supports updating the minor kernel version immediately, which mean is switch_tag=0
func TestAccTencentCloudNeedFixPostgresqlModifySwitchTimePeriodOperationResource_basic(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccStepSetRegion(t, "ap-guangzhou")
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlModifySwitchTimePeriodOperation,
				PreConfig: func() {
					tcacctest.AccStepSetRegion(t, "ap-guangzhou")
					tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON)
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

const testAccPostgresqlModifySwitchTimePeriodOperation = tcacctest.OperationPresetPGSQL + `

resource "tencentcloud_postgresql_modify_switch_time_period_operation" "modify_switch_time_period_operation" {
  db_instance_id = local.pgsql_id
  switch_tag = 0
}

`
