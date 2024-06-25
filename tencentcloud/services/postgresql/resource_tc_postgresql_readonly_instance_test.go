package postgresql_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testPostgresqlReadonlyInstanceResourceKey = "tencentcloud_postgresql_readonly_instance.example"

// go test -i; go test -test.run TestAccTencentCloudPostgresqlReadonlyInstanceResource_basic -v
func TestAccTencentCloudPostgresqlReadonlyInstanceResource_basic(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccStepSetRegion(t, "ap-guangzhou")
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					tcacctest.AccStepSetRegion(t, "ap-guangzhou")
					tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY)
				},
				Config: testAccPostgresqlReadonlyInstanceInstance_basic_without_rogroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "id"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "master_db_instance_id"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "name", "example"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "instance_charge_type", "POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "db_version", "10.23"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "project_id", "0"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "vpc_id"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "subnet_id"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "zone"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "security_groups_ids.#", "1"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "memory", "4"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "cpu", "2"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "storage", "200"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "private_access_ip"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "private_access_port"),
				),
			},
			{
				PreConfig: func() {
					tcacctest.AccStepSetRegion(t, "ap-guangzhou")
					tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY)
				},
				Config: testAccPostgresqlReadonlyInstanceInstance_basic_update_rogroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "id"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "master_db_instance_id"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "name", "example"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "instance_charge_type", "POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "db_version", "10.23"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "project_id", "0"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "vpc_id"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "subnet_id"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "zone"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "security_groups_ids.#", "1"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "memory", "4"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "cpu", "2"),
					resource.TestCheckResourceAttr(testPostgresqlReadonlyInstanceResourceKey, "storage", "200"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "private_access_ip"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "private_access_port"),
					resource.TestCheckResourceAttrSet(testPostgresqlReadonlyInstanceResourceKey, "read_only_group_id"),
				),
			},
		},
	})
}

const testAccPostgresqlReadonlyInstanceInstance_basic_without_rogroup string = `
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
  engine_version    = "10.23"
  root_user         = "root123"
  root_password     = "Root123$"
  charset           = "UTF8"
  project_id        = 0
  memory            = 2
  cpu               = 1
  storage           = 10

  tags = {
    test = "tf"
  }
}

resource "tencentcloud_postgresql_readonly_instance" "example" {
  master_db_instance_id = tencentcloud_postgresql_instance.example.id
  zone                  = var.availability_zone
  name                  = "example"
  auto_renew_flag       = 0
  db_version            = "10.23"
  instance_charge_type  = "POSTPAID_BY_HOUR"
  memory                = 4
  cpu                   = 2
  storage               = 200
  vpc_id                = tencentcloud_vpc.vpc.id
  subnet_id             = tencentcloud_subnet.subnet.id
  need_support_ipv6     = 0
  project_id            = 0
  security_groups_ids   = [
    "sg-5275dorp",
  ]

  tags = {
    createdBy : "terraform"
  }
}
`

const testAccPostgresqlReadonlyInstanceInstance_basic_update_rogroup string = `
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
  engine_version    = "10.23"
  root_user         = "root123"
  root_password     = "Root123$"
  charset           = "UTF8"
  project_id        = 0
  memory            = 2
  cpu               = 1
  storage           = 10

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
  replay_lag_eliminate        = 1
  replay_latency_eliminate    = 1
  max_replay_lag              = 100
  max_replay_latency          = 512
  min_delay_eliminate_reserve = 1
}

resource "tencentcloud_postgresql_readonly_instance" "example" {
  read_only_group_id    = tencentcloud_postgresql_readonly_group.example.id
  master_db_instance_id = tencentcloud_postgresql_instance.example.id
  zone                  = var.availability_zone
  name                  = "example"
  auto_renew_flag       = 0
  db_version            = "10.23"
  instance_charge_type  = "POSTPAID_BY_HOUR"
  memory                = 4
  cpu                   = 2
  storage               = 200
  vpc_id                = tencentcloud_vpc.vpc.id
  subnet_id             = tencentcloud_subnet.subnet.id
  need_support_ipv6     = 0
  project_id            = 0
  security_groups_ids   = [
    "sg-5275dorp",
  ]

  tags = {
    createdBy : "terraform"
  }
}
`
