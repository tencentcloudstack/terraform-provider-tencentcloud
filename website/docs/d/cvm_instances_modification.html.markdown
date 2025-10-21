---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cvm_instances_modification"
sidebar_current: "docs-tencentcloud-datasource-cvm_instances_modification"
description: |-
  Use this data source to query cvm instances modification.
---

# tencentcloud_cvm_instances_modification

Use this data source to query cvm instances modification.

## Example Usage

```hcl
data "tencentcloud_cvm_instances_modification" "foo" {
  instance_ids = ["ins-xxxxxxx"]
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) The upper limit of Filters for each request is 10 and the upper limit for Filter.Values is 2.
* `instance_ids` - (Optional, Set: [`String`]) One or more instance ID to be queried. It can be obtained from the InstanceId in the returned value of API DescribeInstances. The maximum number of instances in batch for each request is 20.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Fields to be filtered.
* `values` - (Required, Set) Value of the field.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_type_config_status_list` - The list of model configurations that can be adjusted by the instance.
  * `instance_type_config` - Configuration information.
    * `cpu` - The number of CPU kernels, in cores.
    * `fpga` - The number of FPGA kernels, in cores.
    * `gpu` - The number of GPU kernels, in cores.
    * `instance_family` - Instance family.
    * `instance_type` - Instance type.
    * `memory` - Memory capacity (in GB).
    * `zone` - Availability zone.
  * `message` - Status description information.
  * `status` - State description.


