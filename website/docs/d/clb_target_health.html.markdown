---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_target_health"
sidebar_current: "docs-tencentcloud-datasource-clb_target_health"
description: |-
  Use this data source to query detailed information of clb target_health
---

# tencentcloud_clb_target_health

Use this data source to query detailed information of clb target_health

## Example Usage

```hcl
data "tencentcloud_clb_target_health" "target_health" {
  load_balancer_ids = ["lb-5dnrkgry"]
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_ids` - (Required, Set: [`String`]) List of IDs of CLB instances to be queried.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `load_balancers` - CLB instance list. Note: This field may return null, indicating that no valid values can be obtained.
  * `listeners` - List of listeners. Note: This field may return null, indicating that no valid values can be obtained.
    * `listener_id` - Listener ID.
    * `listener_name` - Listener name. Note: This field may return null, indicating that no valid values can be obtained.
    * `port` - Listener port.
    * `protocol` - Listener protocol.
    * `rules` - List of forwarding rules of the listener. Note: This field may return null, indicating that no valid values can be obtained.
      * `domain` - Domain name of the forwarding rule. Note: This field may return null, indicating that no valid values can be obtained.
      * `location_id` - Forwarding rule ID.
      * `targets` - Health status of the real server bound to this rule. Note: this field may return null, indicating that no valid values can be obtained.
        * `health_status_detail` - Detailed information about the current health status. Alive: healthy; Dead: exceptional; Unknown: check not started/checking/unknown status.
        * `health_status` - Current health status. true: healthy; false: unhealthy.
        * `ip` - Private IP of the target.
        * `port` - Port bound to the target.
        * `target_id` - Instance ID of the target, such as ins-12345678.
      * `url` - Forwarding rule Url. Note: This field may return null, indicating that no valid values can be obtained.
  * `load_balancer_id` - CLB instance ID.
  * `load_balancer_name` - CLB instance name. Note: This field may return null, indicating that no valid values can be obtained.


