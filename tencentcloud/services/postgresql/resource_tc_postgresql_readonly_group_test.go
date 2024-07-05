package postgresql_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testPostgresqlReadonlyGroupResourceKey = "tencentcloud_postgresql_readonly_group.example"

// go test -i; go test -test.run TestAccTencentCloudPostgresqlReadonlyGroupResource_basic -v
func TestAccTencentCloudPostgresqlReadonlyGroupResource_basic(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccStepSetRegion(t, "ap-guangzhou")
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlReadonlyGroupInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyGroupResourceKey, "id"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyGroupResourceKey, "master_db_instance_id"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyGroupResourceKey, "name", "tf_ro_group"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyGroupResourceKey, "project_id", "0"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyGroupResourceKey, "vpc_id"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyGroupResourceKey, "subnet_id"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyGroupResourceKey, "replay_lag_eliminate", "1"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyGroupResourceKey, "replay_latency_eliminate", "1"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyGroupResourceKey, "max_replay_lag", "100"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyGroupResourceKey, "max_replay_latency", "512"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyGroupResourceKey, "min_delay_eliminate_reserve", "1"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyGroupResourceKey, "net_info_list.#"),
				),
			},
			{
				ResourceName:      testPostgresqlReadonlyGroupResourceKey,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPostgresqlReadonlyGroupInstance_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyGroupResourceKey, "id"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyGroupResourceKey, "master_db_instance_id"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyGroupResourceKey, "name", "tf_ro_group_update"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyGroupResourceKey, "project_id", "0"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyGroupResourceKey, "vpc_id"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyGroupResourceKey, "subnet_id"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyGroupResourceKey, "replay_lag_eliminate", "1"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyGroupResourceKey, "replay_latency_eliminate", "1"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyGroupResourceKey, "max_replay_lag", "100"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyGroupResourceKey, "max_replay_latency", "512"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyGroupResourceKey, "min_delay_eliminate_reserve", "1"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyGroupResourceKey, "net_info_list.#"),
				),
			},
		},
	})
}

const testAccPostgresqlReadonlyGroupInstance string = `
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create postgresql
resource "tencentcloud_postgresql_instance" "example" {
  name              = "example"
  availability_zone = var.availability_zone
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  engine_version    = "10.4"
  root_user         = "root123"
  root_password     = "Root123$"
  charset           = "UTF8"
  project_id        = 0
  memory            = 4
  cpu               = 2
  storage           = 50

  tags = {
    test = "tf"
  }
}

resource "tencentcloud_postgresql_readonly_group" "example" {
  master_db_instance_id       = tencentcloud_postgresql_instance.example.id
  name                        = "tf_ro_group"
  project_id                  = 0
  vpc_id                      = tencentcloud_vpc.vpc.id
  subnet_id                   = tencentcloud_subnet.subnet.id
  security_groups_ids         = ["sg-cm7fbbf3", "sg-5275dorp"]
  replay_lag_eliminate        = 1
  replay_latency_eliminate    = 1
  max_replay_lag              = 100
  max_replay_latency          = 512
  min_delay_eliminate_reserve = 1
}
`

const testAccPostgresqlReadonlyGroupInstance_update string = `
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create postgresql
resource "tencentcloud_postgresql_instance" "example" {
  name              = "example"
  availability_zone = var.availability_zone
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  engine_version    = "10.4"
  root_user         = "root123"
  root_password     = "Root123$"
  charset           = "UTF8"
  project_id        = 0
  memory            = 4
  cpu               = 2
  storage           = 50

  tags = {
    test = "tf"
  }
}

resource "tencentcloud_postgresql_readonly_group" "example" {
  master_db_instance_id       = tencentcloud_postgresql_instance.example.id
  name                        = "tf_ro_group_update"
  project_id                  = 0
  vpc_id                      = tencentcloud_vpc.vpc.id
  subnet_id                   = tencentcloud_subnet.subnet.id
  security_groups_ids         = ["sg-cm7fbbf3"]
  replay_lag_eliminate        = 1
  replay_latency_eliminate    = 1
  max_replay_lag              = 100
  max_replay_latency          = 512
  min_delay_eliminate_reserve = 1
}
`
