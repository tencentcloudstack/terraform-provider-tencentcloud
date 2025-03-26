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
data "tencentcloud_instance_types" "example" {
  availability_zone = "ap-guangzhou-6"
  cpu_core_count    = 4
  memory_size       = 8
}
```

### Complete Example

```hcl
data "tencentcloud_instance_types" "example" {
  cpu_core_count   = 4
  memory_size      = 8
  exclude_sold_out = true

  filter {
    name   = "instance-family"
    values = ["SA2"]
  }

  filter {
    name   = "zone"
    values = ["ap-guangzhou-6"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional, String) The available zone that the CVM instance locates at. This field is conflict with `filter`.
* `cbs_filter` - (Optional, List) Cbs filter.
* `cpu_core_count` - (Optional, Int) The number of CPU cores of the instance.
* `exclude_sold_out` - (Optional, Bool) Indicate to filter instances types that is sold out or not, default is false.
* `filter` - (Optional, Set) One or more name/value pairs to filter. This field is conflict with `availability_zone`.
* `gpu_core_count` - (Optional, Int) The number of GPU cores of the instance.
* `memory_size` - (Optional, Int) Instance memory capacity, unit in GB.
* `result_output_file` - (Optional, String) Used to save results.

The `cbs_filter` object supports the following:

* `disk_charge_type` - (Optional, String) Payment model. Value range:
	- PREPAID: Prepaid;
	- POSTPAID_BY_HOUR: Post-payment.
* `disk_types` - (Optional, List) Hard disk media type. Value range:
	- CLOUD_BASIC: Represents ordinary Cloud Block Storage;
	- CLOUD_PREMIUM: Represents high-performance Cloud Block Storage;
	- CLOUD_SSD: Represents SSD Cloud Block Storage;
	- CLOUD_HSSD: Represents enhanced SSD Cloud Block Storage.
* `disk_usage` - (Optional, String) System disk or data disk. Value range:
	- SYSTEM_DISK: Represents the system disk;
	- DATA_DISK: Represents the data disk.

The `filter` object supports the following:

* `name` - (Required, String) The filter name. Valid values: `zone`, `instance-family` and `instance-charge-type`.
* `values` - (Required, List) The filter values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_types` - An information list of cvm instance. Each element contains the following attributes:
  * `availability_zone` - The available zone that the CVM instance locates at.
  * `cbs_configs` - CBS config. The cbs_configs is populated when the cbs_filter is added.
    * `available` - Whether the configuration is available.
    * `device_class` - Device class.
    * `disk_charge_type` - Payment model. Value range:
	- PREPAID: Prepaid;
	- POSTPAID_BY_HOUR: Post-payment.
    * `disk_type` - Hard disk media type. Value range:
	- CLOUD_BASIC: Represents ordinary Cloud Block Storage;
	- CLOUD_PREMIUM: Represents high-performance Cloud Block Storage;
	- CLOUD_SSD: Represents SSD Cloud Block Storage;
	- CLOUD_HSSD: Represents enhanced SSD Cloud Block Storage.
    * `disk_usage` - Cloud disk type. Value range:
	- SYSTEM_DISK: Represents the system disk;
	- DATA_DISK: Represents the data disk.
    * `extra_performance_range` - Extra performance range.
    * `instance_family` - Instance family.
    * `max_disk_size` - The maximum configurable cloud disk size, in GB.
    * `min_disk_size` - The minimum configurable cloud disk size, in GB.
    * `step_size` - Minimum step size change in cloud disk size, in GB.
    * `zone` - The availability zone to which the Cloud Block Storage belongs.
  * `cpu_core_count` - The number of CPU cores of the instance.
  * `family` - Type series of the instance.
  * `gpu_core_count` - The number of GPU cores of the instance.
  * `instance_charge_type` - Charge type of the instance.
  * `instance_type` - Type of the instance.
  * `memory_size` - Instance memory capacity, unit in GB.
  * `status` - Sell status of the instance.


