---
subcategory: "Anti-DDoS(Dayu)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dayu_ddos_policy_case"
sidebar_current: "docs-tencentcloud-resource-dayu_ddos_policy_case"
description: |-
  Use this resource to create dayu DDoS policy case
---

# tencentcloud_dayu_ddos_policy_case

Use this resource to create dayu DDoS policy case

~> **NOTE:** when a dayu DDoS policy case is created, there will be a dayu DDoS policy created with the same prefix name in the same time. This resource only supports Anti-DDoS of type `bgp`, `bgp-multip` and `bgpip`. One Anti-DDoS resource can only has one DDoS policy case resource. When there is only one Anti-DDoS resource and one policy case, those two resource will be bind automatically.

## Example Usage

```hcl
resource "tencentcloud_dayu_ddos_policy_case" "foo" {
  resource_type       = "bgpip"
  name                = "tf_test_policy_case"
  platform_types      = ["PC", "MOBILE"]
  app_type            = "WEB"
  app_protocols       = ["tcp", "udp"]
  tcp_start_port      = "1000"
  tcp_end_port        = "2000"
  udp_start_port      = "3000"
  udp_end_port        = "4000"
  has_abroad          = "yes"
  has_initiate_tcp    = "yes"
  has_initiate_udp    = "yes"
  peer_tcp_port       = "1111"
  peer_udp_port       = "3333"
  tcp_footprint       = "511"
  udp_footprint       = "500"
  web_api_urls        = ["abc.com", "test.cn/aaa.png"]
  min_tcp_package_len = "1000"
  max_tcp_package_len = "1200"
  min_udp_package_len = "1000"
  max_udp_package_len = "1200"
  has_vpn             = "yes"
}
```

## Argument Reference

The following arguments are supported:

* `app_protocols` - (Required) App protocol set of the DDoS policy case.
* `app_type` - (Required) App type of the DDoS policy case. Valid values: `WEB`, `GAME`, `APP` and `OTHER`.
* `has_abroad` - (Required) Indicate whether the service involves overseas or not. Valid values: `no` and `yes`.
* `has_initiate_tcp` - (Required) Indicate whether the service actively initiates TCP requests or not. Valid values: `no` and `yes`.
* `name` - (Required, ForceNew) Name of the DDoS policy case. Length should between 1 and 64.
* `platform_types` - (Required) Platform set of the DDoS policy case.
* `resource_type` - (Required, ForceNew) Type of the resource that the DDoS policy case works for. Valid values: `bgpip`, `bgp` and `bgp-multip`.
* `tcp_end_port` - (Required) End port of the TCP service. Valid value ranges: (0~65535). It must be greater than `tcp_start_port`.
* `tcp_start_port` - (Required) Start port of the TCP service. Valid value ranges: (0~65535).
* `udp_end_port` - (Required) End port of the UDP service. Valid value ranges: (0~65535). It must be greater than `udp_start_port`.
* `udp_start_port` - (Required) Start port of the UDP service. Valid value ranges: (0~65535).
* `web_api_urls` - (Required) Web API url set.
* `has_initiate_udp` - (Optional) Indicate whether the actively initiate UDP requests or not. Valid values: `no` and `yes`.
* `has_vpn` - (Optional) Indicate whether the service involves VPN service or not. Valid values: `no` and `yes`.
* `max_tcp_package_len` - (Optional) The max length of TCP message package, valid value length should be greater than 0 and less than 1500. It should be greater than `min_tcp_package_len`.
* `max_udp_package_len` - (Optional) The max length of UDP message package, valid value length should be greater than 0 and less than 1500. It should be greater than `min_udp_package_len`.
* `min_tcp_package_len` - (Optional) The minimum length of TCP message package, valid value length should be greater than 0 and less than 1500.
* `min_udp_package_len` - (Optional) The minimum length of UDP message package, valid value length should be greater than 0 and less than 1500.
* `peer_tcp_port` - (Optional) The port that actively initiates TCP requests. Valid value ranges: (1~65535).
* `peer_udp_port` - (Optional) The port that actively initiates UDP requests. Valid value ranges: (1~65535).
* `tcp_footprint` - (Optional) The fixed signature of TCP protocol load, valid value length is range from 1 to 512.
* `udp_footprint` - (Optional) The fixed signature of TCP protocol load, valid value length is range from 1 to 512.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time of the DDoS policy case.
* `scene_id` - Id of the DDoS policy case.


