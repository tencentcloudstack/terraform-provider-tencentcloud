package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testDataPostgresqlInstancesName = "data.tencentcloud_postgresql_instances.id_test"

func TestAccTencentCloudPostgresqlInstancesDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccStepSetRegion(t, "ap-guangzhou")
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataPostgresqlInstanceBasic,
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-guangzhou")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testDataPostgresqlInstancesName, "instance_list.#", "1"),
					resource.TestCheckResourceAttrSet(testDataPostgresqlInstancesName, "instance_list.0.id"),
					resource.TestCheckResourceAttrSet(testDataPostgresqlInstancesName, "instance_list.0.create_time"),
					resource.TestCheckResourceAttrSet(testDataPostgresqlInstancesName, "instance_list.0.charge_type"),
					resource.TestCheckResourceAttrSet(testDataPostgresqlInstancesName, "instance_list.0.engine_version"),
					resource.TestCheckResourceAttrSet(testDataPostgresqlInstancesName, "instance_list.0.project_id"),
					resource.TestCheckResourceAttrSet(testDataPostgresqlInstancesName, "instance_list.0.memory"),
					resource.TestCheckResourceAttrSet(testDataPostgresqlInstancesName, "instance_list.0.storage"),
					resource.TestCheckResourceAttrSet(testDataPostgresqlInstancesName, "instance_list.0.private_access_ip"),
					resource.TestCheckResourceAttrSet(testDataPostgresqlInstancesName, "instance_list.0.private_access_port"),
					resource.TestCheckResourceAttrSet(testDataPostgresqlInstancesName, "instance_list.0.public_access_switch"),
					resource.TestCheckResourceAttrSet(testDataPostgresqlInstancesName, "instance_list.0.charset"),
				),
			},
		},
	})
}

const testAccTencentCloudDataPostgresqlInstanceBasic = CommonPresetPGSQL + `

data "tencentcloud_postgresql_instances" "id_test"{
  id = local.pgsql_id
}
`
