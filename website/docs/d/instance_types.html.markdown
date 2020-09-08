---
subcategory: "CVM"
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
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional) The available zone that the CVM instance locates at. This field is conflict with `filter`.
* `cpu_core_count` - (Optional) The number of CPU cores of the instance.
* `filter` - (Optional) One or more name/value pairs to filter. This field is conflict with `availability_zone`.
* `gpu_core_count` - (Optional) The number of GPU cores of the instance.
* `memory_size` - (Optional) Instance memory capacity, unit in GB.
* `result_output_file` - (Optional) Used to save results.

The `filter` object supports the following:

* `name` - (Required) The filter name, the available values include `zone` and `instance-family`.
* `values` - (Required) The filter values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_types` - An information list of cvm instance. Each element contains the following attributes:
  * `availability_zone` - The available zone that the CVM instance locates at.
  * `cpu_core_count` - The number of CPU cores of the instance.
  * `family` - Type series of the instance.
  * `gpu_core_count` - The number of GPU cores of the instance.
  * `instance_type` - Type of the instance.
  * `memory_size` - Instance memory capacity, unit in GB.


