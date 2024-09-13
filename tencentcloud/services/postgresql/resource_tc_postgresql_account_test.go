package postgresql_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -i; go test -test.run TestAccTencentCloudPostgresqlAccountResource_basic -v
func TestAccTencentCloudPostgresqlAccountResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresAccount,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_postgres_account.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgres_account.example", "db_instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgres_account.example", "user_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgres_account.example", "password"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgres_account.example", "type"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgres_account.example", "remark"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgres_account.example", "lock_status"),
				),
			},
			{
				Config: testAccPostgresAccountUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_postgres_account.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgres_account.example", "db_instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgres_account.example", "user_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgres_account.example", "password"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgres_account.example", "type"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgres_account.example", "remark"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgres_account.example", "lock_status"),
				),
			},
			{
				ResourceName:            "tencentcloud_postgres_account.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

const testAccPostgresAccount = `
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

resource "tencentcloud_postgresql_instance" "example" {
  name              = "example"
  availability_zone = var.availability_zone
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  db_major_version  = "10"
  engine_version    = "10.23"
  root_user         = "root123"
  root_password     = "Root123$"
  charset           = "UTF8"
  project_id        = 0
  cpu               = 1
  memory            = 2
  storage           = 10

  tags = {
    test = "tf"
  }
}

resource "tencentcloud_postgresql_account" "example" {
  db_instance_id = tencentcloud_postgresql_instance.example.id
  user_name      = "tf_example"
  password       = "Password@123"
  type           = "normal"
  remark         = "remark"
  lock_status    = false
}
`

const testAccPostgresAccountUpdate = `
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

resource "tencentcloud_postgresql_instance" "example" {
  name              = "example"
  availability_zone = var.availability_zone
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  db_major_version  = "10"
  engine_version    = "10.23"
  root_user         = "root123"
  root_password     = "Root123$"
  charset           = "UTF8"
  project_id        = 0
  cpu               = 1
  memory            = 2
  storage           = 10

  tags = {
    test = "tf"
  }
}

resource "tencentcloud_postgresql_account" "example" {
  db_instance_id = tencentcloud_postgresql_instance.example.id
  user_name      = "tf_example"
  password       = "Password@456"
  type           = "normal"
  remark         = "remark_update"
  lock_status    = true
}
`
