package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testDataSqlserverInstancesName = "data.tencentcloud_sqlserver_instances.id_test"

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
					testAccCheckSqlserverInstanceExists("tencentcloud_sqlserver_instance.test"),
					resource.TestCheckResourceAttr(testDataSqlserverInstancesName, "instance_list.#", "1"),
					resource.TestCheckResourceAttrSet(testDataSqlserverInstancesName, "instance_list.0.id"),
					resource.TestCheckResourceAttrSet(testDataSqlserverInstancesName, "instance_list.0.create_time"),
					resource.TestCheckResourceAttrSet(testDataSqlserverInstancesName, "instance_list.0.id"),
					resource.TestCheckResourceAttr(testDataSqlserverInstancesName, "instance_list.0.charge_type", "POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttrSet(testDataSqlserverInstancesName, "instance_list.0.engine_version"),
					resource.TestCheckResourceAttr(testDataSqlserverInstancesName, "instance_list.0.project_id", "0"),
					resource.TestCheckResourceAttr(testDataSqlserverInstancesName, "instance_list.0.memory", "2"),
					resource.TestCheckResourceAttr(testDataSqlserverInstancesName, "instance_list.0.storage", "10"),
					resource.TestCheckResourceAttrSet(testDataSqlserverInstancesName, "instance_list.0.vip"),
					resource.TestCheckResourceAttrSet(testDataSqlserverInstancesName, "instance_list.0.vport"),
					resource.TestCheckResourceAttrSet(testDataSqlserverInstancesName, "instance_list.0.status"),
					resource.TestCheckResourceAttrSet(testDataSqlserverInstancesName, "instance_list.0.used_storage"),
				),
			},
		},
	})
}

var testAccTencentCloudDataSqlserverInstancesBasic = testAccSqlserverInstanceBasic + `

resource "tencentcloud_sqlserver_instance" "test" {
  name              = "tf_sqlserver_instance"
  availability_zone = local.az
  charge_type       = "POSTPAID_BY_HOUR"
  project_id        = 0
  memory            = 2
  storage           = 10
}

data "tencentcloud_sqlserver_instances" "id_test"{
	id = tencentcloud_sqlserver_instance.test.id
}
`
