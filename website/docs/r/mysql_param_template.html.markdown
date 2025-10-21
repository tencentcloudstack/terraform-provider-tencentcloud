---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_param_template"
sidebar_current: "docs-tencentcloud-resource-mysql_param_template"
description: |-
  Provides a resource to create a mysql param template
---

# tencentcloud_mysql_param_template

Provides a resource to create a mysql param template

## Example Usage

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

resource "tencentcloud_mysql_param_template" "example" {
  name           = "tf-example"
  description    = "desc."
  engine_version = "8.0"
  param_list {
    current_value = "1"
    name          = "auto_increment_increment"
  }
  param_list {
    current_value = "1"
    name          = "auto_increment_offset"
  }
  param_list {
    current_value = "ON"
    name          = "automatic_sp_privileges"
  }
  template_type = "HIGH_STABILITY"
  engine_type   = "InnoDB"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) The name of parameter template.
* `description` - (Optional, String) The description of parameter template.
* `engine_type` - (Optional, String) The engine type of instance, optional value is InnoDB or RocksDB, default to InnoDB.
* `engine_version` - (Optional, String) The version of MySQL.
* `param_list` - (Optional, List) parameter list.
* `template_id` - (Optional, Int) The ID of source parameter template.
* `template_type` - (Optional, String) The default type of parameter template, supported value is HIGH_STABILITY or HIGH_PERFORMANCE.

The `param_list` object supports the following:

* `current_value` - (Optional, String) The value of parameter.
* `name` - (Optional, String) The name of parameter.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mysql param template can be imported using the id, e.g.

```
terraform import tencentcloud_mysql_param_template.param_template template_id
```

