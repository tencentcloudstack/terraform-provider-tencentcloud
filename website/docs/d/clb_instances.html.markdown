---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_instances"
sidebar_current: "docs-tencentcloud-datasource-clb_instances"
description: |-
  Use this data source to query detailed information of CLB
---

# tencentcloud_clb_instances

Use this data source to query detailed information of CLB

## Example Usage

```hcl
data "tencentcloud_clb_instances" "foo" {
  clb_id             = "lb-k2zjp9lv"
  network_type       = "OPEN"
  clb_name           = "myclb"
  project_id         = "Default Project"
  result_output_file = "mytestpath"
}
```

## Argument Reference

The following arguments are supported:

* `clb_id` - (Optional) ID of the CLB to be queried.
* `clb_name` - (Optional) Name of the CLB to be queried.
* `network_type` - (Optional) Type of CLB instance, and available values include 'OPEN' and 'INTERNAL'
* `project_id` - (Optional) Project ID of the CLB.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `clb_list` - A list of cloud load balancers. Each element contains the following attributes:
  * `clb_id` - ID of CLB.
  * `clb_name` - Name of CLB.
  * `clb_vips` - The virtual service address table of the CLB.
  * `create_time` - Creation time of the CLB
  * `network_type` - Types of CLB.
  * `project_id` - ID of the project.
  * `security_groups` - ID of the security groups.
  * `status_time` - Latest state transition time of CLB.
  * `status` - The status of CLB.
  * `subnet_id` - ID of the subnet
  * `target_region_info_region` - Region information of backend service are attached the CLB.
  * `target_region_info_vpc_id` - VpcId information of backend service are attached the CLB.
  * `vpc_id` - ID of the VPC


