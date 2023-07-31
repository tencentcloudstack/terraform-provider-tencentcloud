---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_renew_db_instance"
sidebar_current: "docs-tencentcloud-resource-sqlserver_renew_db_instance"
description: |-
  Provides a resource to create a sqlserver renew_db_instance
---

# tencentcloud_sqlserver_renew_db_instance

Provides a resource to create a sqlserver renew_db_instance

## Example Usage

```hcl
data "tencentcloud_availability_zones" "zones" {}

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
  charge_type            = "PREPAID"
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
  tags = {
    "createBy" = "tfExample"
  }
}

resource "tencentcloud_sqlserver_renew_db_instance" "renew_db_instance" {
  instance_id = tencentcloud_sqlserver_instance.example.id
  period      = 1
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID.
* `period` - (Optional, Int) How many months to renew, the value range is 1-48, the default is 1.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

sqlserver renew_db_instance can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_renew_db_instance.renew_db_instance renew_db_instance_id
```

