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

### Query with Network and Performance Requirements

```hcl
data "tencentcloud_instance_types" "high_network" {
  availability_zone = "ap-guangzhou-6"
  cpu_core_count    = 8
  memory_size       = 16
}

output "instance_details" {
  value = [for instance in data.tencentcloud_instance_types.high_network.instance_types : {
    type            = instance.instance_type
    type_name       = instance.type_name
    network_card    = instance.network_card
    bandwidth       = instance.instance_bandwidth
    pps             = instance.instance_pps
    cpu_type        = instance.cpu_type
    frequency       = instance.frequency
    status_category = instance.status_category
  }]
}
```

### Query GPU Instances

```hcl
data "tencentcloud_instance_types" "gpu_instances" {
  gpu_core_count = 1

  filter {
    name   = "zone"
    values = ["ap-guangzhou-6"]
  }
}

output "gpu_details" {
  value = [for instance in data.tencentcloud_instance_types.gpu_instances.instance_types : {
    type      = instance.instance_type
    gpu_count = instance.gpu_count
    fpga      = instance.fpga
  }]
}
```

### Query with Local Disk Support

```hcl
data "tencentcloud_instance_types" "local_disk" {
  availability_zone = "ap-guangzhou-6"
  cpu_core_count    = 4
}

output "local_disk_types" {
  value = [for instance in data.tencentcloud_instance_types.local_disk.instance_types :
    instance.local_disk_type_list if length(instance.local_disk_type_list) > 0
  ]
}
```

### Query Price Information

```hcl
data "tencentcloud_instance_types" "with_pricing" {
  availability_zone = "ap-guangzhou-6"
  cpu_core_count    = 2
  memory_size       = 4
}

output "pricing_info" {
  value = [for instance in data.tencentcloud_instance_types.with_pricing.instance_types : {
    type  = instance.instance_type
    price = length(instance.price) > 0 ? instance.price[0] : null
  }]
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
  * `cpu_type` - Processor model.
  * `externals` - Extended attributes.
    * `release_address` - Whether to release address.
    * `storage_block_attr` - HDD local storage attributes.
      * `max_size` - Maximum size of storage block, in GB.
      * `min_size` - Minimum size of storage block, in GB.
      * `type` - Storage block type.
    * `unsupport_networks` - Unsupported network types. Valid values: BASIC (basic network), VPC1.0 (VPC 1.0).
  * `family` - Type series of the instance.
  * `fpga` - Number of FPGA cores.
  * `frequency` - CPU frequency information.
  * `gpu_core_count` - The number of GPU cores of the instance.
  * `gpu_count` - Physical GPU card count mapped to instance. vGPU type is less than 1, direct-attach GPU type is greater than or equal to 1.
  * `instance_bandwidth` - Internal network bandwidth, unit: Gbps.
  * `instance_charge_type` - Charge type of the instance.
  * `instance_pps` - Network packet forwarding capacity, unit: 10K PPS.
  * `instance_type` - Type of the instance.
  * `local_disk_type_list` - List of local disk specifications. Empty if instance type does not support local disks.
    * `max_size` - Maximum size of local disk, in GB.
    * `min_size` - Minimum size of local disk, in GB.
    * `partition_type` - Local disk partition type.
    * `required` - Whether local disk is required when purchasing. Valid values: REQUIRED, OPTIONAL.
    * `type` - Local disk type.
  * `memory_size` - Instance memory capacity, unit in GB.
  * `network_card` - Network card type, for example: 25 represents 25G network card.
  * `price` - Instance pricing information.
    * `charge_unit` - Subsequent billing unit. Valid values: HOUR, GB.
    * `discount_price` - Discount price for prepaid mode, unit: CNY.
    * `discount` - Discount rate. For example, 20.0 means 20% off.
    * `original_price` - Original price for prepaid mode, unit: CNY.
    * `unit_price_discount_second_step` - Subsequent discount unit price for time range (96, 360) hours in postpaid mode, unit: CNY.
    * `unit_price_discount_third_step` - Discounted price of subsequent total cost for usage time interval exceeding 360 hr in postpaid billing mode. measurement unit: usd.
    * `unit_price_discount` - Subsequent discount unit price, used in postpaid mode, unit: CNY.
    * `unit_price_second_step` - Subsequent unit price for time range (96, 360) hours in postpaid mode, unit: CNY.
    * `unit_price_third_step` - Specifies the original price of subsequent total costs with a usage time interval exceeding 360 hr in postpaid billing mode. measurement unit: usd.
    * `unit_price` - Subsequent unit price, used in postpaid mode, unit: CNY.
  * `remark` - Instance remark information.
  * `sold_out_reason` - Reason for sold out status.
  * `status_category` - Stock status category. Valid values: EnoughStock, NormalStock, UnderStock, WithoutStock.
  * `status` - Sell status of the instance.
  * `storage_block_amount` - Number of local storage blocks.
  * `type_name` - Instance type display name.


