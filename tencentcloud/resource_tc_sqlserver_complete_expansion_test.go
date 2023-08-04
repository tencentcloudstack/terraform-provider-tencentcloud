package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverCompleteExpansionResource_basic -v
// go test -v -run TestAccTencentCloudSqlserverCompleteExpansionResource_basic -timeout=0
func TestAccTencentCloudSqlserverCompleteExpansionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: testAccCheckSqlserverInstanceDestroy,
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewSqlserverInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverInstanceExists(testSqlserverInstanceResourceKey),
					resource.TestCheckResourceAttrSet(testSqlserverInstanceResourceKey, "id"),
				),
			},
			{
				Config: testAccUpdateNewSqlserverInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverInstanceExists(testSqlserverInstanceResourceKey),
					resource.TestCheckResourceAttrSet(testSqlserverInstanceResourceKey, "id"),
				),
			},
		},
	})
}

const testAccNewSqlserverInstance string = `
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "example-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  name              = "example-vpc"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_security_group" "security_group" {
  name        = "example-sg"
  description = "desc."
}

resource "tencentcloud_sqlserver_instance" "example" {
  name                   = "tf_example_sql"
  availability_zone      = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type            = "POSTPAID_BY_HOUR"
  vpc_id                 = tencentcloud_vpc.vpc.id
  subnet_id              = tencentcloud_subnet.subnet.id
  security_groups        = [tencentcloud_security_group.security_group.id]
  project_id             = 0
  memory                 = 2
  storage                = 20
  maintenance_week_set   = [1, 2, 3]
  maintenance_start_time = "01:00"
  maintenance_time_span  = 3
  tags                   = {
    "createBy" = "tfExample"
  }
}

resource "tencentcloud_sqlserver_instance" "test" {
  name                   = "tf_sqlserver_instance"
  availability_zone      = "ap-guangzhou-7"
  charge_type            = "POSTPAID_BY_HOUR"
  vpc_id                 = "vpc-1yg5ua6l"
  subnet_id              = "subnet-h7av55g8"
  security_groups        = ["sg-mayqdlt1"]
  project_id             = 0
  memory                 = 2
  storage                = 20
  maintenance_week_set   = [1, 2, 3]
  maintenance_start_time = "09:00"
  maintenance_time_span  = 3
  tags                   = {
    "test" = "test"
  }
}
`

const testAccUpdateNewSqlserverInstance string = `
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "example-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  name              = "example-vpc"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_security_group" "security_group" {
  name        = "example-sg"
  description = "desc."
}

resource "tencentcloud_sqlserver_instance" "example" {
  name                   = "tf_example_sql"
  availability_zone      = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type            = "POSTPAID_BY_HOUR"
  vpc_id                 = tencentcloud_vpc.vpc.id
  subnet_id              = tencentcloud_subnet.subnet.id
  security_groups        = [tencentcloud_security_group.security_group.id]
  project_id             = 0
  memory                 = 4
  storage                = 40
  maintenance_week_set   = [1, 2, 3]
  maintenance_start_time = "01:00"
  maintenance_time_span  = 3
  wait_switch            = 1
  tags                   = {
    "createBy" = "tfExample"
  }
}

resource "tencentcloud_sqlserver_complete_expansion" "example" {
  instance_id = tencentcloud_sqlserver_instance.example.id
}
`
