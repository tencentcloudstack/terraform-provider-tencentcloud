Provides a resource to create a sqlserver migration

Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_security_group" "security_group" {
  name        = "sg-example"
  description = "desc."
}

resource "tencentcloud_sqlserver_basic_instance" "src_example" {
  name                   = "tf-example"
  availability_zone      = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type            = "POSTPAID_BY_HOUR"
  vpc_id                 = tencentcloud_vpc.vpc.id
  subnet_id              = tencentcloud_subnet.subnet.id
  project_id             = 0
  memory                 = 4
  storage                = 100
  cpu                    = 2
  machine_type           = "CLOUD_PREMIUM"
  maintenance_week_set   = [1, 2, 3]
  maintenance_start_time = "09:00"
  maintenance_time_span  = 3
  security_groups        = [tencentcloud_security_group.security_group.id]

  tags = {
    "test" = "test"
  }
}

resource "tencentcloud_sqlserver_basic_instance" "dst_example" {
  name                   = "tf-example"
  availability_zone      = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type            = "POSTPAID_BY_HOUR"
  vpc_id                 = tencentcloud_vpc.vpc.id
  subnet_id              = tencentcloud_subnet.subnet.id
  project_id             = 0
  memory                 = 4
  storage                = 100
  cpu                    = 2
  machine_type           = "CLOUD_PREMIUM"
  maintenance_week_set   = [1, 2, 3]
  maintenance_start_time = "09:00"
  maintenance_time_span  = 3
  security_groups        = [tencentcloud_security_group.security_group.id]

  tags = {
    "test" = "test"
  }
}

resource "tencentcloud_sqlserver_db" "src" {
  instance_id = tencentcloud_sqlserver_basic_instance.src_example.id
  name        = "tf_example_db_src"
  charset     = "Chinese_PRC_BIN"
  remark      = "testACC-remark"
}

resource "tencentcloud_sqlserver_db" "dst" {
  instance_id = tencentcloud_sqlserver_basic_instance.dst_example.id
  name        = "tf_example_db_dst"
  charset     = "Chinese_PRC_BIN"
  remark      = "testACC-remark"
}

resource "tencentcloud_sqlserver_account" "src" {
  instance_id = tencentcloud_sqlserver_basic_instance.src_example.id
  name        = "tf_example_src_account"
  password    = "Qwer@234"
  is_admin    = true
}

resource "tencentcloud_sqlserver_account" "dst" {
  instance_id = tencentcloud_sqlserver_basic_instance.dst_example.id
  name        = "tf_example_dst_account"
  password    = "Qwer@234"
  is_admin    = true
}

resource "tencentcloud_sqlserver_account_db_attachment" "src" {
  instance_id  = tencentcloud_sqlserver_basic_instance.src_example.id
  account_name = tencentcloud_sqlserver_account.src.name
  db_name      = tencentcloud_sqlserver_db.src.name
  privilege    = "ReadWrite"
}

resource "tencentcloud_sqlserver_account_db_attachment" "dst" {
  instance_id  = tencentcloud_sqlserver_basic_instance.dst_example.id
  account_name = tencentcloud_sqlserver_account.dst.name
  db_name      = tencentcloud_sqlserver_db.dst.name
  privilege    = "ReadWrite"
}

resource "tencentcloud_sqlserver_migration" "migration" {
  migrate_name = "tf_test_migration"
  migrate_type = 1
  source_type  = 1
  source {
    instance_id = tencentcloud_sqlserver_basic_instance.src_example.id
    user_name   = tencentcloud_sqlserver_account.src.name
    password    = tencentcloud_sqlserver_account.src.password
  }
  target {
    instance_id = tencentcloud_sqlserver_basic_instance.dst_example.id
    user_name   = tencentcloud_sqlserver_account.dst.name
    password    = tencentcloud_sqlserver_account.dst.password
  }

  migrate_db_set {
    db_name = tencentcloud_sqlserver_db.src.name
  }
}
```

Import

sqlserver migration can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_migration.migration migration_id
```