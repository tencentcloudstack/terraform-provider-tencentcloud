---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_instance_types"
sidebar_current: "docs-tencentcloud-datasource-instance_types"
description: |-
  Use this data source to query instances type.
---

# tencentcloud_instance_types

Use this data source to query instances type.

## Example Usage

```hcl
data "tencentcloud_instance_types" "foo" {
  availability_zone = "ap-guangzhou-2"
  cpu_core_count    = 2
  memory_size       = 4
}

data tencentcloud_instance_types "t1c1g" {
  cpu_core_count   = 1
  memory_size      = 1
  exclude_sold_out = true
  filter {
    name   = "instance-charge-type"
    values = ["POSTPAID_BY_HOUR"]
  }
  filter {
    name   = "zone"
    values = ["ap-shanghai-2"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional, String) The available zone that the CVM instance locates at. This field is conflict with `filter`.
* `cpu_core_count` - (Optional, Int) The number of CPU cores of the instance.
* `exclude_sold_out` - (Optional, Bool) Indicate to filter instances types that is sold out or not, default is false.
* `filter` - (Optional, Set) One or more name/value pairs to filter. This field is conflict with `availability_zone`.
* `gpu_core_count` - (Optional, Int) The number of GPU cores of the instance.
* `memory_size` - (Optional, Int) Instance memory capacity, unit in GB.
* `result_output_file` - (Optional, String) Used to save results.

The `filter` object supports the following:

* `name` - (Required, String) The filter name. Valid values: `zone`, `instance-family` and `instance-charge-type`.
* `values` - (Required, List) The filter values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_types` - An information list of cvm instance. Each element contains the following attributes:
  * `availability_zone` - The available zone that the CVM instance locates at.
  * `cpu_core_count` - The number of CPU cores of the instance.
  * `family` - Type series of the instance.
  * `gpu_core_count` - The number of GPU cores of the instance.
  * `instance_charge_type` - Charge type of the instance.
  * `instance_type` - Type of the instance.
  * `memory_size` - Instance memory capacity, unit in GB.
  * `status` - Sell status of the instance.


