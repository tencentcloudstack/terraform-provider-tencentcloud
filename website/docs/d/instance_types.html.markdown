---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_instance_types"
sidebar_current: "docs-tencentcloud-datasource-instance-types"
description: |-
  Provides a list of instance types available to the user.
---

# tencentcloud_instance_types

The Instance Types data source list the cvm_instance_types of TencentCloud.

## Example Usage

```hcl
data "tencentcloud_instance_types" "lowest_cost_config" {
  filter {
    name   = "instance-family"
    values = ["S1"]
  }

  cpu_core_count = 1
  memory_size    = 1
}
```

## Argument Reference

* `filter` - (Optional) One or more name/value pairs to filter off of. There are several valid keys:  `zone`,`instance-family`. For a full reference, check out [DescribeInstanceTypeConfigs in the TencentCloud API reference](https://intl.cloud.tencent.com/document/api/213/9391).
 * `cpu_core_count` - (Optional) Limit search to specific cpu core count.
 * `memory_size` -  (Optional) Limit search to specific memory size.

## Attributes Reference

The following attributes are exported

 * `availability_zone` - Indicate the availability zone for this instance type.
 * `instance_type` - TencentCloud instance type of the cvm instance.
 * `cpu_core_count` - Number of CPU cores.
 * `memory_size` - Size of memory, measured in GB.
 * `family` - The instance type family.