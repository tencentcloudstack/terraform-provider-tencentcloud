Provides a resource to create a sqlserver config_backup_strategy

Example Usage

Daily backup

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

resource "tencentcloud_sqlserver_basic_instance" "example" {
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

resource "tencentcloud_sqlserver_config_backup_strategy" "example" {
  instance_id              = tencentcloud_sqlserver_basic_instance.example.id
  backup_type              = "daily"
  backup_time              = 0
  backup_day               = 1
  backup_model             = "master_no_pkg"
  backup_cycle             = [1]
  backup_save_days         = 7
  regular_backup_enable    = "disable"
  regular_backup_save_days = 90
  regular_backup_strategy  = "months"
  regular_backup_counts    = 1
}
```

Weekly backup

```hcl
resource "tencentcloud_sqlserver_config_backup_strategy" "example" {
  instance_id              = tencentcloud_sqlserver_basic_instance.example.id
  backup_type              = "weekly"
  backup_time              = 0
  backup_model             = "master_no_pkg"
  backup_cycle             = [1, 3, 5]
  backup_save_days         = 7
  regular_backup_enable    = "disable"
  regular_backup_save_days = 90
  regular_backup_strategy  = "months"
  regular_backup_counts    = 1
}
```

Regular backup

```hcl
resource "tencentcloud_sqlserver_config_backup_strategy" "example" {
  instance_id               = tencentcloud_sqlserver_basic_instance.example.id
  backup_time               = 0
  backup_model              = "master_no_pkg"
  backup_cycle              = [1, 3]
  backup_save_days          = 7
  regular_backup_enable     = "enable"
  regular_backup_save_days  = 120
  regular_backup_strategy   = "months"
  regular_backup_counts     = 1
  regular_backup_start_time = "%s"
}
```

Import

sqlserver config_backup_strategy can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_config_backup_strategy.example mssql-si2823jyl
```