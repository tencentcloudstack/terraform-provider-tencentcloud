package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccPostgresqlrecoverytimeObject = "data.tencentcloud_postgresql_recovery_time.recovery_time"

func TestAccTencentCloudPostgresqlRecoveryTimeDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccStepSetRegion(t, "ap-guangzhou")
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlRecoveryTimeDataSource,
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-guangzhou")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(testAccPostgresqlrecoverytimeObject),
					resource.TestCheckResourceAttrSet(testAccPostgresqlrecoverytimeObject, "db_instance_id"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlrecoverytimeObject, "recovery_begin_time"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlrecoverytimeObject, "recovery_end_time"),
				),
			},
		},
	})
}

const testAccPostgresqlRecoveryTimeDataSource = CommonPresetPGSQL + `

data "tencentcloud_postgresql_recovery_time" "recovery_time" {
  db_instance_id = local.pgsql_id
}

`
