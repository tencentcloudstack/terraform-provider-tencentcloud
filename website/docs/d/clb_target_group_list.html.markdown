---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_target_group_list"
sidebar_current: "docs-tencentcloud-datasource-clb_target_group_list"
description: |-
  Use this data source to query detailed information of clb target_group_list
---

# tencentcloud_clb_target_group_list

Use this data source to query detailed information of clb target_group_list

## Example Usage

```hcl
data "tencentcloud_clb_target_group_list" "target_group_list" {
  filters {
    name   = "TargetGroupName"
    values = ["keep-tgg"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter array, which is exclusive of TargetGroupIds. Valid values: TargetGroupVpcId and TargetGroupName. Target group ID will be used first.
* `result_output_file` - (Optional, String) Used to save results.
* `target_group_ids` - (Optional, Set: [`String`]) Target group ID array.

The `filters` object supports the following:

* `name` - (Required, String) Filter name.
* `values` - (Required, Set) Filter value array.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `target_group_set` - Information set of displayed target groups.
  * `associated_rule` - Array of associated rules. Note: this field may return null, indicating that no valid values can be obtained.
    * `domain` - Domain name of associated forwarding rule. Note: this field may return null, indicating that no valid values can be obtained.
    * `listener_id` - ID of associated listener.
    * `listener_name` - Listener name.
    * `load_balancer_id` - ID of associated CLB instance.
    * `load_balancer_name` - CLB instance name.
    * `location_id` - ID of associated forwarding rule. Note: this field may return null, indicating that no valid values can be obtained.
    * `port` - Port of associated listener.
    * `protocol` - Protocol type of associated listener, such as HTTP or TCP.
    * `url` - URL of associated forwarding rule. Note: this field may return null, indicating that no valid values can be obtained.
  * `created_time` - Target group creation time.
  * `port` - Default port of target group. Note: this field may return null, indicating that no valid values can be obtained.
  * `target_group_id` - Target group ID.
  * `target_group_name` - Target group name.
  * `updated_time` - Target group modification time.
  * `vpc_id` - vpcid of target group.


