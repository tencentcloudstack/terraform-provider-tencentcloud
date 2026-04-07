---
subcategory: "Vector Database(VDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vdb_instance"
sidebar_current: "docs-tencentcloud-resource-vdb_instance"
description: |-
  Provides a resource to create a VDB instance.
---

# tencentcloud_vdb_instance

Provides a resource to create a VDB instance.

~> **NOTE:** When `force_delete` is false (default), destroying this resource will only isolate the instance to the recycle bin. Set `force_delete` to true to permanently destroy the instance.

## Example Usage

### Create a pay-as-you-go single instance

```hcl
resource "tencentcloud_vdb_instance" "example" {
  vpc_id             = "vpc-xxxxxxxx"
  subnet_id          = "subnet-xxxxxxxx"
  pay_mode           = 0
  security_group_ids = ["sg-xxxxxxxx"]
  instance_name      = "tf-example"
  instance_type      = "single"
  node_type          = "normal"
  cpu                = 2
  memory             = 4
  disk_size          = 50
  worker_node_num    = 1
  force_delete       = false
}
```

### Create a monthly subscription cluster instance with all parameters

```hcl
resource "tencentcloud_vdb_instance" "cluster" {
  vpc_id          = "vpc-xxxxxxxx"
  subnet_id       = "subnet-xxxxxxxx"
  pay_mode        = 1
  pay_period      = 1
  auto_renew      = 1
  instance_name   = "tf-example-cluster"
  instance_type   = "cluster"
  mode            = "two"
  product_type    = 0
  node_type       = "compute"
  cpu             = 4
  memory          = 8
  disk_size       = 100
  worker_node_num = 2
  params          = "{\"key\":\"value\"}"
  force_delete    = true

  security_group_ids = ["sg-xxxxxxxx"]

  resource_tags {
    tag_key   = "env"
    tag_value = "test"
  }

  resource_tags {
    tag_key   = "project"
    tag_value = "demo"
  }
}
```

## Argument Reference

The following arguments are supported:

* `pay_mode` - (Required, Int) Billing mode. 0: pay-as-you-go, 1: monthly subscription.
* `security_group_ids` - (Required, List: [`String`]) Security group IDs.
* `subnet_id` - (Required, String, ForceNew) Subnet ID.
* `vpc_id` - (Required, String, ForceNew) VPC ID.
* `auto_renew` - (Optional, Int) Auto-renew flag. 0: disabled, 1: enabled.
* `cpu` - (Optional, Int) CPU cores.
* `disk_size` - (Optional, Int) Disk size in GB.
* `force_delete` - (Optional, Bool) Whether to force delete (destroy) the instance. If false, only isolate to recycle bin. If true, isolate then destroy. Default is false.
* `instance_name` - (Optional, String) Instance name. Supports up to 60 characters.
* `instance_type` - (Optional, String, ForceNew) Instance type. Valid values: base, single, cluster.
* `memory` - (Optional, Int) Memory size in GB.
* `mode` - (Optional, String, ForceNew) Availability zone mode for cluster type. Valid values: two, three.
* `node_type` - (Optional, String, ForceNew) Node type. Valid values: compute, normal, store.
* `params` - (Optional, String) Instance extra parameters, submitted via JSON.
* `pay_period` - (Optional, Int) Monthly subscription period in months. Valid values: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36. Default is 1.
* `product_type` - (Optional, Int, ForceNew) Product version. 0: standard, 1: capacity-enhanced.
* `resource_tags` - (Optional, List) Tag list.
* `worker_node_num` - (Optional, Int) Number of worker nodes.

The `resource_tags` object supports the following:

* `tag_key` - (Required, String) Tag key.
* `tag_value` - (Required, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `api_version` - API version.
* `created_at` - Creation time.
* `engine_name` - Engine name.
* `engine_version` - Engine version.
* `expired_at` - Expiration time.
* `extend` - Extended information in JSON format.
* `is_no_expired` - Whether the instance never expires.
* `isolate_at` - Isolation time.
* `networks` - Network information.
  * `expire_time` - Old IP expiration time.
  * `port` - Internal port.
  * `preserve_duration` - Old IP preservation duration in days.
  * `subnet_id` - Subnet ID.
  * `vip` - Internal IP.
  * `vpc_id` - VPC ID.
* `nodes` - Instance node list.
  * `name` - Pod name.
  * `status` - Pod status.
* `product` - Product.
* `region` - Region.
* `shard_num` - Shard number.
* `status` - Instance status.
* `task_status` - Task status. 0: no task, 1: pending, 2-11: various operations in progress.
* `wan_address` - Public network address.
* `zone` - Availability zone.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `30m`) Used when creating the resource.
* `update` - (Defaults to `30m`) Used when updating the resource.
* `delete` - (Defaults to `10m`) Used when deleting the resource.

## Import

VDB instance can be imported using the id, e.g.

```
terraform import tencentcloud_vdb_instance.example vdb-xxxxxxxx
```

