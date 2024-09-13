package postgresql_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -i; go test -test.run TestAccTencentCloudPostgresqlAccountPrivilegesOperationResource_basic -v
func TestAccTencentCloudPostgresqlAccountPrivilegesOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresAccountPrivileges,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_account_privileges_operation.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_account_privileges_operation.example", "db_instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_account_privileges_operation.example", "user_name"),
				),
			},
		},
	})
}

const testAccPostgresAccountPrivileges = `
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

# create account
resource "tencentcloud_postgresql_account" "example" {
  db_instance_id = tencentcloud_postgresql_instance.example.id
  user_name      = "tf_example"
  password       = "Password@123"
  type           = "normal"
  remark         = "remark"
  lock_status    = false
}

# create account privileges
resource "tencentcloud_postgresql_account_privileges_operation" "example" {
  db_instance_id = tencentcloud_postgresql_instance.example.id
  user_name      = tencentcloud_postgresql_account.example.user_name
  modify_privilege_set {
    database_privilege {
      object {
        object_name = "postgres"
        object_type = "database"
      }

      privilege_set = ["CONNECT", "TEMPORARY", "CREATE"]
    }

    modify_type = "grantObject"
    is_cascade  = false
  }
}
`
