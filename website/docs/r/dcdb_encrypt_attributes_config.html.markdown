---
subcategory: "TDSQL for MySQL(DCDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dcdb_encrypt_attributes_config"
sidebar_current: "docs-tencentcloud-resource-dcdb_encrypt_attributes_config"
description: |-
  Provides a resource to create a dcdb encrypt_attributes_config
---

# tencentcloud_dcdb_encrypt_attributes_config

Provides a resource to create a dcdb encrypt_attributes_config

~> **NOTE:**  This resource currently only supports the newly created MySQL 8.0.24 version.

## Example Usage

```hcl
data "tencentcloud_security_groups" "internal" {
  name = "default"
}

data "tencentcloud_vpc_instances" "vpc" {
  name = "Default-VPC"
}

data "tencentcloud_vpc_subnets" "subnet" {
  vpc_id = data.tencentcloud_vpc_instances.vpc.instance_list.0.vpc_id
}

locals {
  vpc_id    = data.tencentcloud_vpc_subnets.subnet.instance_list.0.vpc_id
  subnet_id = data.tencentcloud_vpc_subnets.subnet.instance_list.0.subnet_id
  sg_id     = data.tencentcloud_security_groups.internal.security_groups.0.security_group_id
}

resource "tencentcloud_dcdb_db_instance" "prepaid_instance" {
  instance_name    = "test_dcdb_db_post_instance"
  zones            = [var.default_az]
  period           = 1
  shard_memory     = "2"
  shard_storage    = "10"
  shard_node_count = "2"
  shard_count      = "2"
  vpc_id           = local.vpc_id
  subnet_id        = local.subnet_id
  db_version_id    = "8.0"
  resource_tags {
    tag_key   = "aaa"
    tag_value = "bbb"
  }
  security_group_ids = [local.sg_id]
}

resource "tencentcloud_dcdb_hourdb_instance" "hourdb_instance" {
  instance_name     = "test_dcdb_db_hourdb_instance"
  zones             = [var.default_az]
  shard_memory      = "2"
  shard_storage     = "10"
  shard_node_count  = "2"
  shard_count       = "2"
  vpc_id            = local.vpc_id
  subnet_id         = local.subnet_id
  security_group_id = local.sg_id
  db_version_id     = "8.0"
  resource_tags {
    tag_key   = "aaa"
    tag_value = "bbb"
  }
}

locals {
  prepaid_dcdb_id = tencentcloud_dcdb_db_instance.prepaid_instance.id
  hourdb_dcdb_id  = tencentcloud_dcdb_hourdb_instance.hourdb_instance.id
}

// for postpaid instance
resource "tencentcloud_dcdb_encrypt_attributes_config" "config_hourdb" {
  instance_id     = local.hourdb_dcdb_id
  encrypt_enabled = 1
}

// for prepaid instance
resource "tencentcloud_dcdb_encrypt_attributes_config" "config_prepaid" {
  instance_id     = local.prepaid_dcdb_id
  encrypt_enabled = 1
}
```

## Argument Reference

The following arguments are supported:

* `encrypt_enabled` - (Required, Int) whether to enable data encryption. Notice: it is not supported to turn it off after it is turned on. The optional values: 0-disable, 1-enable.
* `instance_id` - (Required, String) instance id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

dcdb encrypt_attributes_config can be imported using the id, e.g.

```
terraform import tencentcloud_dcdb_encrypt_attributes_config.encrypt_attributes_config encrypt_attributes_config_id
```

