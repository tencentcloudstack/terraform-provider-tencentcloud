Provides a mysql account privilege resource to grant different access privilege to different database. A database can be granted by multiple account.

~> **NOTE:** It has been deprecated and replaced by  tencentcloud_mysql_privilege.

Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "cdb"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-mysql"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  name              = "subnet-mysql"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_security_group" "security_group" {
  name        = "sg-mysql"
  description = "mysql test"
}

resource "tencentcloud_mysql_instance" "example" {
  internet_service  = 1
  engine_version    = "5.7"
  charge_type       = "POSTPAID"
  root_password     = "PassWord123"
  slave_deploy_mode = 0
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  slave_sync_mode   = 1
  instance_name     = "tf-example-mysql"
  mem_size          = 4000
  volume_size       = 200
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  intranet_port     = 3306
  security_groups   = [tencentcloud_security_group.security_group.id]

  tags = {
    name = "test"
  }

  parameters = {
    character_set_server = "utf8"
    max_connections      = "1000"
  }
}

resource "tencentcloud_mysql_account" "example" {
  mysql_id             = tencentcloud_mysql_instance.example.id
  name                 = "tf_example"
  password             = "Qwer@234"
  description          = "desc."
  max_user_connections = 10
}

resource "tencentcloud_mysql_account_privilege" "default" {
  mysql_id       = tencentcloud_mysql_instance.example.id
  account_name   = tencentcloud_mysql_account.example.name
  privileges     = ["SELECT", "INSERT", "UPDATE", "DELETE"]
  database_names = ["dbname1", "dbname2"]
}
```