---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dayu_ddos_policy_cases"
sidebar_current: "docs-tencentcloud-datasource-dayu_ddos_policy_cases"
description: |-
  Use this data source to query dayu DDoS policy cases
---

# tencentcloud_dayu_ddos_policy_cases

Use this data source to query dayu DDoS policy cases

## Example Usage

```hcl
data "tencentcloud_dayu_ddos_policy_cases" "id_test" {
  resource_type = tencentcloud_dayu_ddos_policy_case.test_policy_case.resource_type
  scene_id      = tencentcloud_dayu_ddos_policy_case.test_policy_case.scene_id
}
```

## Argument Reference

The following arguments are supported:

* `resource_type` - (Required) Type of the resource that the DDoS policy case works for, valid values are `bgpip`, `bgp`, `bgp-multip`, `net`.
* `scene_id` - (Required) Id of the DDoS policy case to be query.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of DDoS policy cases. Each element contains the following attributes.
  * `app_protocols` - App protocol set of the DDoS policy case.
  * `app_type` - App type of the DDoS policy case.
  * `create_time` - Create time of the DDoS policy case.
  * `has_abroad` - Indicate whether the service involves overseas or not.
  * `has_initiate_tcp` - Indicate whether the service actively initiates TCP requests or not.
  * `has_initiate_udp` - Indicate whether the actively initiate UDP requests or not.
  * `has_vpn` - Indicate whether the service involves VPN service or not.
  * `max_tcp_package_len` - The max length of TCP message package.
  * `max_udp_package_len` - The max length of UDP message package.
  * `min_tcp_package_len` - The minimum length of TCP message package.
  * `min_udp_package_len` - The minimum length of UDP message package.
  * `name` - Name of the DDoS policy case.
  * `peer_tcp_port` - The port that actively initiates TCP requests.
  * `peer_udp_port` - The port that actively initiates UDP requests.
  * `platform_types` - Platform set of the DDoS policy case.
  * `resource_type` - Type of the resource that the DDoS policy case works for.
  * `scene_id` - Id of the DDoS policy case.
  * `tcp_end_port` - End port of the TCP service.
  * `tcp_foot_print` - The fixed signature of TCP protocol load.
  * `tcp_start_port` - Start port of the TCP service.
  * `udp_end_port` - End port of the UDP service.
  * `udp_foot_print` - The fixed signature of TCP protocol load.
  * `udp_start_port` - Start port of the UDP service.
  * `web_api_urls` - Web API url set.


