Use this data source to query the list of SQL Server account DB privileges.

Example Usage

```hcl
data "tencentcloud_availability_zones" "zones" {}

data "tencentcloud_sqlserver_account_db_attachments" "test" {
  instance_id  = tencentcloud_sqlserver_instance.example.id
  account_name = tencentcloud_sqlserver_account_db_attachment.example.account_name
}

resource "tencentcloud_vpc" "vpc" {
  name       = "example-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones.zones.zones.0.name
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
  availability_zone      = data.tencentcloud_availability_zones.zones.zones.0.name
  charge_type            = "POSTPAID_BY_HOUR"
  period                 = 1
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

resource "tencentcloud_sqlserver_db" "example" {
  instance_id = tencentcloud_sqlserver_instance.example.id
  name        = "tfExampleDb"
  charset     = "Chinese_PRC_BIN"
  remark      = "remark desc."
}

resource "tencentcloud_sqlserver_account" "example" {
  instance_id = tencentcloud_sqlserver_instance.example.id
  name        = "tf_example_account"
  password    = "PassWord@123"
  remark      = "remark desc."
}

resource "tencentcloud_sqlserver_account_db_attachment" "example" {
  instance_id  = tencentcloud_sqlserver_instance.example.id
  account_name = tencentcloud_sqlserver_account.example.name
  db_name      = tencentcloud_sqlserver_db.example.name
  privilege    = "ReadWrite"
}
```