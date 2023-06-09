---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_idle_instances"
sidebar_current: "docs-tencentcloud-datasource-clb_idle_instances"
description: |-
  Use this data source to query detailed information of clb idle_loadbalancers
---

# tencentcloud_clb_idle_instances

Use this data source to query detailed information of clb idle_loadbalancers

## Example Usage

```hcl
data "tencentcloud_clb_idle_instances" "idle_instance" {
  load_balancer_region = "ap-guangzhou"
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_region` - (Optional, String) CLB instance region.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `idle_load_balancers` - List of idle CLBs. Note: This field may return null, indicating that no valid values can be obtained.
  * `domain` - The load balancing hostname.Note: This field may return null, indicating that no valid values can be obtained.
  * `forward` - CLB type. Value range: 1 (CLB); 0 (classic CLB).
  * `idle_reason` - The reason why the load balancer is considered idle. NO_RULES: No rules configured. NO_RS: The rules are not associated with servers.
  * `load_balancer_id` - CLB instance ID.
  * `load_balancer_name` - CLB instance name.
  * `region` - CLB instance region.
  * `status` - CLB instance status, including:0: Creating; 1: Running.
  * `vip` - CLB instance VIP.


