---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_net_detect_states"
sidebar_current: "docs-tencentcloud-datasource-vpc_net_detect_states"
description: |-
  Use this data source to query detailed information of vpc net_detect_states
---

# tencentcloud_vpc_net_detect_states

Use this data source to query detailed information of vpc net_detect_states

## Example Usage

```hcl
data "tencentcloud_vpc_net_detect_states" "net_detect_states" {
  net_detect_ids = ["netd-12345678"]
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter conditions. `NetDetectIds` and `Filters` cannot be specified at the same time.net-detect-id - String - (Filter condition) The network detection instance ID, such as netd-12345678.
* `net_detect_ids` - (Optional, Set: [`String`]) The array of network detection instance `IDs`, such as [`netd-12345678`].
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) The attribute name. If more than one Filter exists, the logical relation between these Filters is `AND`.
* `values` - (Required, Set) Attribute value. If multiple values exist in one filter, the logical relationship between these values is `OR`. For a `bool` parameter, the valid values include `TRUE` and `FALSE`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `net_detect_state_set` - The array of network detection verification results that meet requirements.Note: This field may return null, indicating that no valid values can be obtained.
  * `net_detect_id` - The ID of a network detection instance, such as netd-12345678.
  * `net_detect_ip_state_set` - The array of network detection destination IP verification results.
    * `delay` - The latency. Unit: ms.
    * `detect_destination_ip` - The destination IPv4 address of network detection.
    * `packet_loss_rate` - The packet loss rate.
    * `state` - The detection result.0: successful;-1: no packet loss occurred during routing;-2: packet loss occurred when outbound traffic is blocked by the ACL;-3: packet loss occurred when inbound traffic is blocked by the ACL;-4: other errors.


