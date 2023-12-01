---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_cluster_slave_zone"
sidebar_current: "docs-tencentcloud-resource-cynosdb_cluster_slave_zone"
description: |-
  Provides a resource to create a cynosdb cluster slave zone.
---

# tencentcloud_cynosdb_cluster_slave_zone

Provides a resource to create a cynosdb cluster slave zone.

## Example Usage

### Set a new slave zone for a cynosdb cluster.

```hcl
locals {
  vpc_id    = data.tencentcloud_vpc_subnets.gz3.instance_list.0.vpc_id
  subnet_id = data.tencentcloud_vpc_subnets.gz3.instance_list.0.subnet_id
  sg_id     = data.tencentcloud_security_groups.internal.security_groups.0.security_group_id
  sg_id2    = data.tencentcloud_security_groups.exclusive.security_groups.0.security_group_id
}

variable "fixed_tags" {
  default = {
    fixed_resource : "do_not_remove"
  }
}

variable "availability_zone" {
  default = "ap-guangzhou-4"
}

variable "new_availability_zone" {
  default = "ap-guangzhou-6"
}

variable "my_param_template" {
  default = "15765"
}

data "tencentcloud_security_groups" "internal" {
  name = "default"
  tags = var.fixed_tags
}

data "tencentcloud_security_groups" "exclusive" {
  name = "test_preset_sg"
}

data "tencentcloud_vpc_subnets" "gz3" {
  availability_zone = var.default_az
  is_default        = true
}

resource "tencentcloud_cynosdb_cluster" "instance" {
  available_zone               = var.availability_zone
  vpc_id                       = local.vpc_id
  subnet_id                    = local.subnet_id
  db_type                      = "MYSQL"
  db_version                   = "5.7"
  storage_limit                = 1000
  cluster_name                 = "tf_test_cynosdb_cluster_slave_zone"
  password                     = "cynos@123"
  instance_maintain_duration   = 3600
  instance_maintain_start_time = 10800
  instance_maintain_weekdays = [
    "Fri",
    "Mon",
    "Sat",
    "Sun",
    "Thu",
    "Wed",
    "Tue",
  ]

  instance_cpu_core    = 1
  instance_memory_size = 2
  param_items {
    name          = "character_set_server"
    current_value = "utf8"
  }
  param_items {
    name          = "time_zone"
    current_value = "+09:00"
  }

  force_delete = true

  rw_group_sg = [
    local.sg_id
  ]
  ro_group_sg = [
    local.sg_id
  ]
  prarm_template_id = var.my_param_template
}

resource "tencentcloud_cynosdb_cluster_slave_zone" "cluster_slave_zone" {
  cluster_id = tencentcloud_cynosdb_cluster.instance.id
  slave_zone = var.new_availability_zone
}
```

### Update the slave zone with specified value.

```hcl
resource "tencentcloud_cynosdb_cluster_slave_zone" "cluster_slave_zone" {
  cluster_id = tencentcloud_cynosdb_cluster.instance.id
  slave_zone = var.availability_zone

  timeouts {
    create = "500s"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) The ID of cluster.
* `slave_zone` - (Required, String) Slave zone.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cynosdb cluster_slave_zone can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_cluster_slave_zone.cluster_slave_zone cluster_id#slave_zone
```

