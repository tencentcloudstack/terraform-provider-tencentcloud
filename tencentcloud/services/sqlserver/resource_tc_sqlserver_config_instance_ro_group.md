Provides a resource to create a sqlserver config_instance_ro_group

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

resource "tencentcloud_sqlserver_readonly_instance" "example" {
  name                     = "tf_example"
  availability_zone        = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type              = "POSTPAID_BY_HOUR"
  vpc_id                   = tencentcloud_vpc.vpc.id
  subnet_id                = tencentcloud_subnet.subnet.id
  memory                   = 4
  storage                  = 20
  master_instance_id       = tencentcloud_sqlserver_basic_instance.example.id
  readonly_group_type      = 2
  read_only_group_name     = "tf_example_ro"
  is_offline_delay         = 1
  read_only_max_delay_time = 10
  min_read_only_in_group   = 0
  force_upgrade            = true
}

resource "tencentcloud_sqlserver_config_instance_ro_group" "example" {
  instance_id              = tencentcloud_sqlserver_readonly_instance.example.master_instance_id
  read_only_group_id       = tencentcloud_sqlserver_readonly_instance.example.readonly_group_id
  read_only_group_name     = "tf_example_ro_update"
  is_offline_delay         = 1
  read_only_max_delay_time = 5
  min_read_only_in_group   = 1
}
```

Import

sqlserver config_instance_ro_group can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_config_instance_ro_group.example mssql-ds1xhnt9#mssqlro-o6dv2ugx#0#0
```