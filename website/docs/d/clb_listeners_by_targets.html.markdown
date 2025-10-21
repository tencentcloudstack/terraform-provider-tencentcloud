---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_listeners_by_targets"
sidebar_current: "docs-tencentcloud-datasource-clb_listeners_by_targets"
description: |-
  Use this data source to query detailed information of clb listeners_by_targets
---

# tencentcloud_clb_listeners_by_targets

Use this data source to query detailed information of clb listeners_by_targets

## Example Usage

```hcl
data "tencentcloud_clb_listeners_by_targets" "listeners_by_targets" {
  backends {
    vpc_id     = "vpc-4owdpnwr"
    private_ip = "106.52.160.211"
  }
}
```

## Argument Reference

The following arguments are supported:

* `backends` - (Required, List) List of private network IPs to be queried.
* `result_output_file` - (Optional, String) Used to save results.

The `backends` object supports the following:

* `private_ip` - (Required, String) Private network IP to be queried, which can be of the CVM or ENI.
* `vpc_id` - (Required, String) VPC ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `load_balancers` - Detail of the CLB instance.
  * `listeners` - Listener rule.
    * `end_port` - End port of the listener. Note: this field may return null, indicating that no valid values can be obtained.
    * `listener_id` - Listener ID.
    * `port` - Listener port.
    * `protocol` - Listener protocol.
    * `rules` - Bound rule. Note: this field may return null, indicating that no valid values can be obtained.
      * `domain` - Domain name.
      * `location_id` - Rule ID.
      * `targets` - Object bound to the real server.
        * `port` - Port bound to the real server.
        * `private_ip` - Private network IP of the real server.
        * `type` - Private network IP type, which can be cvm or eni.
        * `vpc_id` - VPC ID of the real server. Note: this field may return null, indicating that no valid values can be obtained.
        * `weight` - Weight of the real server. Note: this field may return null, indicating that no valid values can be obtained.
      * `url` - url.
    * `targets` - Object bound to the layer-4 listener. Note: this field may return null, indicating that no valid values can be obtained.
      * `port` - Port bound to the real server.
      * `private_ip` - Private network IP of the real server.
      * `type` - Private network IP type, which can be cvm or eni.
      * `vpc_id` - VPC ID of the real server. Note: this field may return null, indicating that no valid values can be obtained.
      * `weight` - Weight of the real server. Note: this field may return null, indicating that no valid values can be obtained.
  * `load_balancer_id` - String ID of the CLB instance.
  * `region` - CLB instance region.
  * `vip` - VIP of the CLB instance.


