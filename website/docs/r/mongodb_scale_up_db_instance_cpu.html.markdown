---
subcategory: "TencentDB for MongoDB(mongodb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mongodb_scale_up_db_instance_cpu"
sidebar_current: "docs-tencentcloud-resource-mongodb_scale_up_db_instance_cpu"
description: |-
  Provides a resource to manage MongoDB instance elastic CPU scale-up.
---

# tencentcloud_mongodb_scale_up_db_instance_cpu

Provides a resource to manage MongoDB instance elastic CPU scale-up.

~> **NOTE:** This resource is used to manually enable elastic CPU scale-up for MongoDB instances. When the resource is created, `ScaleUpDBInstanceCpu` is called to scale up the CPU. When the resource is destroyed, `ScaleDownDBInstanceCpu` is called to disable the elastic CPU scale-up and shrink back.

## Example Usage

```hcl
resource "tencentcloud_mongodb_instance" "example" {
  instance_name  = "tf-example"
  cpu            = 2
  memory         = 4
  volume         = 100
  engine_version = "MONGO_80_WT"
  machine_type   = "GE.LD.T1"
  available_zone = "ap-guangzhou-6"
  vpc_id         = "vpc-i5yyodl9"
  subnet_id      = "subnet-hhi88a58"
  project_id     = 0
  password       = "Password@2026"
}

resource "tencentcloud_mongodb_scale_up_db_instance_cpu" "example" {
  instance_id = tencentcloud_mongodb_instance.example.id
  extra_cpu   = 2
}
```

## Argument Reference

The following arguments are supported:

* `extra_cpu` - (Required, Int, ForceNew) Number of CPU cores to scale up (increment). Each node will increase by this number of cores. It cannot be changed after creation.
* `instance_id` - (Required, String, ForceNew) MongoDB instance ID, for example: cmgo-p8vnipr5.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `is_locked` - Whether the instance is locked (a flow is in progress).
* `is_scaled` - Whether the instance is currently in scaled state.
* `max_extra_cpu` - Maximum number of CPU cores that can be scaled up, MIN(minimum shard CPU cores, 24).
* `reason` - Reason why elastic CPU scale-up is not available.
* `scalable` - Whether elastic CPU scale-up is available.
* `scale_up_time` - Scale-up time.
* `trigger_type` - Scale-up trigger type. Valid values: 1-manual, 2-periodic, 3-one-time time period, 4-monitor.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `30m`) Used when creating the resource.
* `delete` - (Defaults to `30m`) Used when deleting the resource.

## Import

MongoDB instance elastic CPU scale-up can be imported using the instance id, e.g.

```
terraform import tencentcloud_mongodb_scale_up_db_instance_cpu.example cmgo-p8vnipr5
```

