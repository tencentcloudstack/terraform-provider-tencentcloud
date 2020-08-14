package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testDataSqlserverInstancesName = "data.tencentcloud_sqlserver_instances.id_test"

func TestAccTencentCloudDataSqlserverInstances(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBDestroy,
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

const testAccTencentCloudDataSqlserverInstancesBasic = `
variable "availability_zone"{
default = "ap-guangzhou-2"
}

resource "tencentcloud_sqlserver_instance" "test" {
  name = "tf_postsql_instance"
  availability_zone = var.availability_zone
  charge_type = "POSTPAID_BY_HOUR"
  project_id = 0
  memory = 2
  storage = 10
}

data "tencentcloud_sqlserver_instances" "id_test"{
	id = tencentcloud_sqlserver_instance.test.id
}
`
