package cdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlRoGroupResource_basic -v -timeout=0
func TestAccTencentCloudMysqlRoGroupResource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlRoGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_ro_group.ro_group", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_ro_group.ro_group", "ro_group_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_ro_group.ro_group", "ro_group_info.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_ro_group.ro_group", "ro_group_info.0.min_ro_in_group"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_ro_group.ro_group", "ro_group_info.0.replication_delay_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_ro_group.ro_group", "ro_group_info.0.ro_group_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_ro_group.ro_group", "ro_group_info.0.ro_max_delay_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_ro_group.ro_group", "ro_group_info.0.ro_offline_delay"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_ro_group.ro_group", "ro_group_info.0.weight_mode"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_ro_group.ro_group", "ro_weight_values.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_ro_group.ro_group", "ro_weight_values.0.instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_ro_group.ro_group", "ro_weight_values.0.weight"),
				),
			},
		},
	})
}

const testAccMysqlRoGroup = `
variable "security_groups" {
	default = "` + tcacctest.DefaultCrsSecurityGroups + `"
}

variable "availability_zone" {
  default = "ap-guangzhou-4"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-mysql"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "subnet-mysql"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_mysql_instance" "example" {
  internet_service  = 1
  engine_version    = "5.7"
  charge_type       = "POSTPAID"
  root_password     = "PassWord123"
  slave_deploy_mode = 0
  availability_zone = var.availability_zone
  slave_sync_mode   = 1
  instance_name     = "tf-example-mysql"
  mem_size          = 2000
  volume_size       = 200
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  intranet_port     = 3306
  security_groups   = [var.security_groups]

  tags = {
    name = "test"
  }

  parameters = {
    character_set_server = "utf8"
    max_connections      = "1000"
  }
}

resource "tencentcloud_mysql_readonly_instance" "foo" {
  master_instance_id = tencentcloud_mysql_instance.example.id
  instance_name      = "tf-mysql"
  mem_size           = 2000
  volume_size        = 200
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_id          = tencentcloud_subnet.subnet.id
  intranet_port      = 3306
  security_groups    = [var.security_groups]

  tags = {
    createBy = "terraform"
  }
}

data "tencentcloud_mysql_instance" "foo" {
  mysql_id = tencentcloud_mysql_instance.example.id

  depends_on = [tencentcloud_mysql_readonly_instance.foo]
}

resource "tencentcloud_mysql_ro_group" "ro_group" {
	instance_id = tencentcloud_mysql_instance.example.id
	ro_group_id = data.tencentcloud_mysql_instance.foo.instance_list.0.ro_groups.0.group_id
	ro_group_info {
	  ro_group_name          = "keep-ro"
	  ro_max_delay_time      = 1
	  ro_offline_delay       = 1
	  min_ro_in_group        = 1
	  weight_mode            = "custom"
	  # replication_delay_time = 1
	}
	ro_weight_values {
	  instance_id = tencentcloud_mysql_readonly_instance.foo.id
	  weight      = 10
	}
	is_balance_ro_load = 1
}
`
