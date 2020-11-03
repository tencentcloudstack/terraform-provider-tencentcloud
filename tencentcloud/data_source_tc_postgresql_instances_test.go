package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testDataPostgresqlInstancesName = "data.tencentcloud_postgresql_instances.id_test"

func TestAccTencentCloudDataPostgresqlInstances(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataPostgresqlInstanceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPostgresqlInstanceExists("tencentcloud_postgresql_instance.test"),
					resource.TestCheckResourceAttr(testDataPostgresqlInstancesName, "instance_list.#", "1"),
					resource.TestCheckResourceAttrSet(testDataPostgresqlInstancesName, "instance_list.0.id"),
					resource.TestCheckResourceAttrSet(testDataPostgresqlInstancesName, "instance_list.0.create_time"),
					resource.TestCheckResourceAttr(testDataPostgresqlInstancesName, "instance_list.0.charge_type", "POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttr(testDataPostgresqlInstancesName, "instance_list.0.engine_version", "9.3.5"),
					resource.TestCheckResourceAttr(testDataPostgresqlInstancesName, "instance_list.0.project_id", "0"),
					resource.TestCheckResourceAttr(testDataPostgresqlInstancesName, "instance_list.0.memory", "2"),
					resource.TestCheckResourceAttr(testDataPostgresqlInstancesName, "instance_list.0.storage", "10"),
					resource.TestCheckResourceAttrSet(testDataPostgresqlInstancesName, "instance_list.0.private_access_ip"),
					resource.TestCheckResourceAttrSet(testDataPostgresqlInstancesName, "instance_list.0.private_access_port"),
					resource.TestCheckResourceAttrSet(testDataPostgresqlInstancesName, "instance_list.0.public_access_switch"),
					resource.TestCheckResourceAttrSet(testDataPostgresqlInstancesName, "instance_list.0.charset"),
					resource.TestCheckResourceAttr(testDataPostgresqlInstancesName, "instance_list.0.tags.tf", "test"),
				),
			},
		},
	})
}

const testAccTencentCloudDataPostgresqlInstanceBasic = `
variable "availability_zone"{
default = "ap-guangzhou-2"
}

resource "tencentcloud_postgresql_instance" "test" {
name = "tf_postsql_instance"
availability_zone = var.availability_zone
charge_type = "POSTPAID_BY_HOUR"
engine_version		= "9.3.5"
root_password                 = "1qaA2k1wgvfa3ZZZ"
charset = "UTF8"
project_id = 0
memory = 2
storage = 10

	tags = {
		tf = "test"
	}
}

data "tencentcloud_postgresql_instances" "id_test"{
	id = tencentcloud_postgresql_instance.test.id
}
`
