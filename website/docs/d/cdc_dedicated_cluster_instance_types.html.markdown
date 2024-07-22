---
subcategory: "CDC"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cdc_dedicated_cluster_instance_types"
sidebar_current: "docs-tencentcloud-datasource-cdc_dedicated_cluster_instance_types"
description: |-
  Use this data source to query detailed information of CDC dedicated cluster instance types
---

# tencentcloud_cdc_dedicated_cluster_instance_types

Use this data source to query detailed information of CDC dedicated cluster instance types

## Example Usage

```hcl
data "tencentcloud_cdc_dedicated_cluster_instance_types" "types" {
  dedicated_cluster_id = "cluster-262n63e8"
}
```

## Argument Reference

The following arguments are supported:

* `dedicated_cluster_id` - (Required, String) Dedicated Cluster ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `dedicated_cluster_instance_type_set` - Dedicated Cluster Supported InstanceType.
  * `cpu_type` - Instance CPU Type.
  * `cpu` - Instance CPU.
  * `fpga` - Instance Fpga.
  * `gpu` - Instance GPU.
  * `instance_bandwidth` - Instance Bandwidth.
  * `instance_family` - Instance Family.
  * `instance_pps` - Instance Pps.
  * `instance_type` - Instance Type.
  * `memory` - Instance Memory.
  * `network_card` - Instance Type.
  * `remark` - Instance Remark.
  * `status` - Instance Status.
  * `storage_block_amount` - Instance Storage Block Amount.
  * `type_name` - Instance Type Name.
  * `zone` - Zone Name.


