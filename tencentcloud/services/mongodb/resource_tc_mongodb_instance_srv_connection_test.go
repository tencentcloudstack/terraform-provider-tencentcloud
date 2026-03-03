package mongodb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMongodbInstanceSrvConnectionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbInstanceSrvConnectionBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance_srv_connection.srv", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance_srv_connection.srv", "domain"),
				),
			},
			{
				ResourceName:      "tencentcloud_mongodb_instance_srv_connection.srv",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudMongodbInstanceSrvConnectionResource_customDomain(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbInstanceSrvConnectionCustomDomain,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance_srv_connection.srv", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance_srv_connection.srv", "domain", "test.mongodb.com"),
				),
			},
			{
				Config: testAccMongodbInstanceSrvConnectionCustomDomainUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance_srv_connection.srv", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance_srv_connection.srv", "domain", "updated.mongodb.com"),
				),
			},
		},
	})
}

const testAccMongodbInstanceSrvConnectionBasic = tcacctest.DefaultMongoDBSpec + `
resource "tencentcloud_vpc" "vpc" {
	name       = "mongodb-srv-vpc"
	cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	name              = "mongodb-srv-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_mongodb_instance" "mongodb" {
	instance_name  = "tf-mongodb-srv-test"
	memory         = local.memory
	volume         = local.volume
	engine_version = local.engine_version
	machine_type   = local.machine_type
	security_groups = [local.security_group_id]
	available_zone = "ap-guangzhou-3"
	project_id     = 0
	password       = "test1234"
	vpc_id         = tencentcloud_vpc.vpc.id
	subnet_id      = tencentcloud_subnet.subnet.id
}

resource "tencentcloud_mongodb_instance_srv_connection" "srv" {
	instance_id = tencentcloud_mongodb_instance.mongodb.id
}
`

const testAccMongodbInstanceSrvConnectionCustomDomain = tcacctest.DefaultMongoDBSpec + `
resource "tencentcloud_vpc" "vpc" {
	name       = "mongodb-srv-vpc"
	cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	name              = "mongodb-srv-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_mongodb_instance" "mongodb" {
	instance_name  = "tf-mongodb-srv-test-custom"
	memory         = local.memory
	volume         = local.volume
	engine_version = local.engine_version
	machine_type   = local.machine_type
	security_groups = [local.security_group_id]
	available_zone = "ap-guangzhou-3"
	project_id     = 0
	password       = "test1234"
	vpc_id         = tencentcloud_vpc.vpc.id
	subnet_id      = tencentcloud_subnet.subnet.id
}

resource "tencentcloud_mongodb_instance_srv_connection" "srv" {
	instance_id = tencentcloud_mongodb_instance.mongodb.id
	domain      = "test.mongodb.com"
}
`

const testAccMongodbInstanceSrvConnectionCustomDomainUpdate = tcacctest.DefaultMongoDBSpec + `
resource "tencentcloud_vpc" "vpc" {
	name       = "mongodb-srv-vpc"
	cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	name              = "mongodb-srv-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_mongodb_instance" "mongodb" {
	instance_name  = "tf-mongodb-srv-test-custom"
	memory         = local.memory
	volume         = local.volume
	engine_version = local.engine_version
	machine_type   = local.machine_type
	security_groups = [local.security_group_id]
	available_zone = "ap-guangzhou-3"
	project_id     = 0
	password       = "test1234"
	vpc_id         = tencentcloud_vpc.vpc.id
	subnet_id      = tencentcloud_subnet.subnet.id
}

resource "tencentcloud_mongodb_instance_srv_connection" "srv" {
	instance_id = tencentcloud_mongodb_instance.mongodb.id
	domain      = "updated.mongodb.com"
}
`
