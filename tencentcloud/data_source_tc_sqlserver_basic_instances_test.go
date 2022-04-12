package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testDataSqlserverBasicInstancesName = "data.tencentcloud_sqlserver_basic_instances.id_test"

func TestAccDataSourceTencentCloudSqlserverBasicInstances(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSqlserverBasicInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataSqlserverBasicInstancesBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverBasicInstanceExists("tencentcloud_sqlserver_basic_instance.test"),
					resource.TestCheckResourceAttr(testDataSqlserverBasicInstancesName, "instance_list.#", "1"),
					resource.TestCheckResourceAttrSet(testDataSqlserverBasicInstancesName, "instance_list.0.id"),
					resource.TestCheckResourceAttrSet(testDataSqlserverBasicInstancesName, "instance_list.0.create_time"),
					resource.TestCheckResourceAttrSet(testDataSqlserverBasicInstancesName, "instance_list.0.name"),
					resource.TestCheckResourceAttr(testDataSqlserverBasicInstancesName, "instance_list.0.charge_type", "POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttrSet(testDataSqlserverBasicInstancesName, "instance_list.0.engine_version"),
					resource.TestCheckResourceAttr(testDataSqlserverBasicInstancesName, "instance_list.0.project_id", "0"),
					resource.TestCheckResourceAttr(testDataSqlserverBasicInstancesName, "instance_list.0.memory", "4"),
					resource.TestCheckResourceAttr(testDataSqlserverBasicInstancesName, "instance_list.0.storage", "20"),
					resource.TestCheckResourceAttr(testDataSqlserverBasicInstancesName, "instance_list.0.cpu", "2"),
					resource.TestCheckResourceAttrSet(testDataSqlserverBasicInstancesName, "instance_list.0.vpc_id"),
					resource.TestCheckResourceAttrSet(testDataSqlserverBasicInstancesName, "instance_list.0.subnet_id"),
					resource.TestCheckResourceAttrSet(testDataSqlserverBasicInstancesName, "instance_list.0.availability_zone"),
					resource.TestCheckResourceAttrSet(testDataSqlserverBasicInstancesName, "instance_list.0.vip"),
					resource.TestCheckResourceAttrSet(testDataSqlserverBasicInstancesName, "instance_list.0.vport"),
					resource.TestCheckResourceAttrSet(testDataSqlserverBasicInstancesName, "instance_list.0.status"),
					resource.TestCheckResourceAttrSet(testDataSqlserverBasicInstancesName, "instance_list.0.used_storage"),
					resource.TestCheckResourceAttr(testDataSqlserverBasicInstancesName, "instance_list.0.tags.test", "test"),
				),
			},
		},
	})
}

const testAccTencentCloudDataSqlserverBasicInstancesBasic = testAccSqlserverInstanceBasic + `

resource "tencentcloud_vpc" "foo" {
	name       = "tf-sqlserver-vpc"
	cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "foo" {
	availability_zone = local.az1
	name              = "tf-sqlserver-subnet"
	vpc_id            = tencentcloud_vpc.foo.id
	cidr_block        = "10.0.0.0/16"
	is_multicast      = false
}

resource "tencentcloud_sqlserver_basic_instance" "test" {
	name                    = "tf_sqlserver_basic_instance"
	availability_zone       = local.az1
	charge_type             = "POSTPAID_BY_HOUR"
	vpc_id                  = tencentcloud_vpc.foo.id
	subnet_id               = tencentcloud_subnet.foo.id
	machine_type            = "CLOUD_PREMIUM"
	project_id              = 0
	memory                  = 4
	storage                 = 20
	cpu                     = 2
	security_groups         = ["` + defaultSecurityGroup + `"]

	tags = {
		"test" = "test"
	}
}

data "tencentcloud_sqlserver_basic_instances" "id_test"{
	id = tencentcloud_sqlserver_basic_instance.test.id
}
`
