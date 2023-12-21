package postgresql_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccPostgresqlrecoverytimeObject = "data.tencentcloud_postgresql_recovery_time.recovery_time"

func TestAccTencentCloudPostgresqlRecoveryTimeDataSource_basic(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccStepSetRegion(t, "ap-guangzhou")
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlRecoveryTimeDataSource,
				PreConfig: func() {
					tcacctest.AccStepSetRegion(t, "ap-guangzhou")
					tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(testAccPostgresqlrecoverytimeObject),
					resource.TestCheckResourceAttrSet(testAccPostgresqlrecoverytimeObject, "db_instance_id"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlrecoverytimeObject, "recovery_begin_time"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlrecoverytimeObject, "recovery_end_time"),
				),
			},
		},
	})
}

const testAccPostgresqlRecoveryTimeDataSource = tcacctest.CommonPresetPGSQL + `

data "tencentcloud_postgresql_recovery_time" "recovery_time" {
  db_instance_id = local.pgsql_id
}

`
