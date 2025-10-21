---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_cross_targets"
sidebar_current: "docs-tencentcloud-datasource-clb_cross_targets"
description: |-
  Use this data source to query detailed information of clb cross_targets
---

# tencentcloud_clb_cross_targets

Use this data source to query detailed information of clb cross_targets

## Example Usage

```hcl
data "tencentcloud_clb_cross_targets" "cross_targets" {
  filters {
    name   = "vpc-id"
    values = ["vpc-4owdpnwr"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter conditions to query CVMs and ENIs: vpc-id - String - Required: No - (Filter condition) Filter by VPC ID, such as vpc-12345678. ip - String - Required: No - (Filter condition) Filter by real server IP, such as 192.168.0.1. listener-id - String - Required: No - (Filter condition) Filter by listener ID, such as lbl-12345678. location-id - String - Required: No - (Filter condition) Filter by forwarding rule ID of the layer-7 listener, such as loc-12345678.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Filter name.
* `values` - (Required, Set) Filter values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cross_target_set` - Cross target set.
  * `eni_id` - ENI ID of the CVM instance.
  * `instance_id` - ID of the CVM instance.Note: This field may return null, indicating that no valid value was found.
  * `instance_name` - Name of the CVM instance. Note: This field may return null, indicating that no valid value was found.
  * `ip` - IP address of the CVM or ENI instance.
  * `local_vpc_id` - VPC ID of the CLB instance.
  * `region` - Region of the CVM or ENI instance.
  * `vpc_id` - VPC ID of the CVM or ENI instance.
  * `vpc_name` - VPC name of the CVM or ENI instance.


