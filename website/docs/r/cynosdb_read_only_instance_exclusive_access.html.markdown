---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_read_only_instance_exclusive_access"
sidebar_current: "docs-tencentcloud-resource-cynosdb_read_only_instance_exclusive_access"
description: |-
  Provides a resource to create a cynosdb read_only_instance_exclusive_access
---

# tencentcloud_cynosdb_read_only_instance_exclusive_access

Provides a resource to create a cynosdb read_only_instance_exclusive_access

## Example Usage

```hcl
variable "cynosdb_cluster_id" {
  default = "default_cynosdb_cluster"
}
variable "cynosdb_cluster_instance_id" {
  default = "default_cluster_instance"
}

variable "cynosdb_cluster_security_group_id" {
  default = "default_security_group_id"
}

data "tencentcloud_vpc_subnets" "gz3" {
  availability_zone = var.default_az
  is_default        = true
}

locals {
  vpc_id    = data.tencentcloud_vpc_subnets.gz3.instance_list.0.vpc_id
  subnet_id = data.tencentcloud_vpc_subnets.gz3.instance_list.0.subnet_id
}

resource "tencentcloud_cynosdb_read_only_instance_exclusive_access" "read_only_instance_exclusive_access" {
  cluster_id         = var.cynosdb_cluster_id
  instance_id        = var.cynosdb_cluster_instance_id
  vpc_id             = local.vpc_id
  subnet_id          = local.subnet_id
  port               = 1234
  security_group_ids = [var.cynosdb_cluster_security_group_id]
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) Cluster ID.
* `instance_id` - (Required, String, ForceNew) Need to activate a read-only instance ID with unique access.
* `port` - (Required, Int, ForceNew) port.
* `subnet_id` - (Required, String, ForceNew) The specified subnet ID.
* `vpc_id` - (Required, String, ForceNew) Specified VPC ID.
* `security_group_ids` - (Optional, Set: [`String`], ForceNew) Security Group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



