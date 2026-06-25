---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cdb_start_cpu_expand"
sidebar_current: "docs-tencentcloud-resource-cdb_start_cpu_expand"
description: |-
  Provides a resource to create a CDB CPU elastic expand attachment
---

# tencentcloud_cdb_start_cpu_expand

Provides a resource to create a CDB CPU elastic expand attachment

## Example Usage

### Auto expansion type

```hcl
resource "tencentcloud_cdb_start_cpu_expand" "example" {
  instance_id = "cdb-test1234"
  type        = "auto"

  auto_strategy {
    expand_threshold     = 80
    shrink_threshold     = 20
    expand_second_period = 300
    shrink_second_period = 600
  }
}
```

### Manual expansion type

```hcl
resource "tencentcloud_cdb_start_cpu_expand" "example" {
  instance_id = "cdb-test1234"
  type        = "manual"
  expand_cpu  = 4
}
```

### TimeInterval expansion type

```hcl
resource "tencentcloud_cdb_start_cpu_expand" "example" {
  instance_id = "cdb-test1234"
  type        = "timeInterval"
  expand_cpu  = 4

  time_interval_strategy {
    start_time = 1709251200
    end_time   = 1709337600
  }
}
```

### Period expansion type

```hcl
resource "tencentcloud_cdb_start_cpu_expand" "example" {
  instance_id = "cdb-test1234"
  type        = "period"
  expand_cpu  = 4

  period_strategy {
    time_cycle {
      monday    = true
      tuesday   = true
      wednesday = true
      thursday  = true
      friday    = true
      saturday  = false
      sunday    = false
    }

    time_interval {
      start_time = "09:00"
      end_time   = "18:00"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID, which can be obtained from the DescribeDBInstances API.
* `type` - (Required, String, ForceNew) Expansion type. Valid values: `auto`, `manual`, `timeInterval`, `period`.
* `auto_strategy` - (Optional, List, ForceNew) Auto expansion strategy. Required when `type` is `auto`.
* `expand_cpu` - (Optional, Int, ForceNew) CPU cores to expand. Required when `type` is `manual`, `timeInterval`, or `period`.
* `period_strategy` - (Optional, List, ForceNew) Period expansion strategy. Required when `type` is `period`.
* `time_interval_strategy` - (Optional, List, ForceNew) Time interval expansion strategy. Required when `type` is `timeInterval`.

The `auto_strategy` object supports the following:

* `expand_threshold` - (Required, Int) Auto expansion threshold. Valid values: 40, 50, 60, 70, 80, 90.
* `shrink_threshold` - (Required, Int) Auto shrink threshold. Valid values: 10, 20, 30.
* `expand_second_period` - (Optional, Int) Expansion observation period in seconds. Valid values: 15, 30, 45, 60, 180, 300, 600, 900, 1800.
* `shrink_second_period` - (Optional, Int) Shrink observation period in seconds. Valid values: 300, 600, 900, 1800.

The `period_strategy` object supports the following:

* `time_cycle` - (Optional, List) Weekly cycle configuration.
* `time_interval` - (Optional, List) Daily time range configuration.

The `time_cycle` object of `period_strategy` supports the following:

* `friday` - (Optional, Bool) Whether to expand on Friday.
* `monday` - (Optional, Bool) Whether to expand on Monday.
* `saturday` - (Optional, Bool) Whether to expand on Saturday.
* `sunday` - (Optional, Bool) Whether to expand on Sunday.
* `thursday` - (Optional, Bool) Whether to expand on Thursday.
* `tuesday` - (Optional, Bool) Whether to expand on Tuesday.
* `wednesday` - (Optional, Bool) Whether to expand on Wednesday.

The `time_interval_strategy` object supports the following:

* `end_time` - (Required, Int) End expansion time as integer timestamp in seconds.
* `start_time` - (Required, Int) Start expansion time as integer timestamp in seconds.

The `time_interval` object of `period_strategy` supports the following:

* `end_time` - (Optional, String) End time string.
* `start_time` - (Optional, String) Start time string.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `async_request_id` - Async request ID returned by Create/Delete APIs.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `20m`) Used when creating the resource.
* `delete` - (Defaults to `20m`) Used when deleting the resource.

## Import

CDB start cpu expand can be imported using the instance_id, e.g.

```
terraform import tencentcloud_cdb_start_cpu_expand.example cdb-test1234
```

