package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testDataSqlserverInstancesName = "data.tencentcloud_sqlserver_instances.test"

func TestAccDataSourceTencentCloudSqlserverInstances(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSqlserverInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataSqlserverInstancesBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testDataSqlserverInstancesName, "instance_list.#"),
					resource.TestCheckResourceAttrSet(testDataSqlserverInstancesName, "instance_list.0.id"),
					resource.TestCheckResourceAttrSet(testDataSqlserverInstancesName, "instance_list.0.create_time"),
					resource.TestCheckResourceAttrSet(testDataSqlserverInstancesName, "instance_list.0.id"),
					resource.TestCheckResourceAttr(testDataSqlserverInstancesName, "instance_list.0.charge_type", "POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttrSet(testDataSqlserverInstancesName, "instance_list.0.engine_version"),
					resource.TestCheckResourceAttrSet(testDataSqlserverInstancesName, "instance_list.0.memory"),
					resource.TestCheckResourceAttrSet(testDataSqlserverInstancesName, "instance_list.0.storage"),
					resource.TestCheckResourceAttrSet(testDataSqlserverInstancesName, "instance_list.0.vip"),
					resource.TestCheckResourceAttrSet(testDataSqlserverInstancesName, "instance_list.0.vport"),
					resource.TestCheckResourceAttrSet(testDataSqlserverInstancesName, "instance_list.0.status"),
					resource.TestCheckResourceAttrSet(testDataSqlserverInstancesName, "instance_list.0.used_storage"),
				),
			},
		},
	})
}

var testAccTencentCloudDataSqlserverInstancesBasic = `

data "tencentcloud_sqlserver_instances" "test"{
  name = "keep"
}
`
