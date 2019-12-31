---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dayu_ddos_policies"
sidebar_current: "docs-tencentcloud-datasource-dayu_ddos_policies"
description: |-
  Use this data source to query dayu DDoS policies
---

# tencentcloud_dayu_ddos_policies

Use this data source to query dayu DDoS policies

## Example Usage

```hcl
data "tencentcloud_dayu_ddos_policies" "id_test" {
  resource_type = tencentcloud_dayu_ddos_policy.test_policy.resource_type
  policy_id     = tencentcloud_dayu_ddos_policy.test_policy.policy_id
}
```

## Argument Reference

The following arguments are supported:

* `resource_type` - (Required) Type of the resource that the DDoS policy works for, valid values are `bgpip`, `bgp`, `bgp-multip`, `net`.
* `policy_id` - (Optional) Id of the DDoS policy to be query.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of DDoS policies. Each element contains the following attributes.
  * `black_white_ips` - Black and white ip list.
    * `ip` - Ip.
    * `type` - Type of the ip.
  * `create_time` - Create time of the DDoS policy.
  * `drop_options` - Option list of abnormal check of the DDoS policy.
    * `bad_conn_threshold` - The number of new connections based on destination IP that trigger suppression of connections.
    * `check_sync_conn` - Indicate whether to check null connection or not.
    * `conn_timeout` - Connection timeout of abnormal connection check.
    * `drop_icmp` - Indicate whether to drop ICMP protocol or not.
    * `drop_other` - Indicate whether to drop other protocols(exclude TCP/UDP/ICMP) or not.
    * `drop_tcp` - Indicate whether to drop TCP protocol or not.
    * `drop_udp` - Indicate to drop UDP protocol or not.
    * `dst_conn_limit` - The limit of concurrent connections based on destination IP.
    * `dst_new_limit` - The limit of new connections based on destination IP.
    * `icmp_mbps_limit` - The limit of ICMP traffic rate.
    * `null_conn_enable` - Indicate to enable null connection or not.
    * `other_mbps_limit` - The limit of other protocols(exclude TCP/UDP/ICMP) traffic rate.
    * `source_conn_limit` - The limit of concurrent connections based on source IP.
    * `source_new_limit` - The limit of new connections based on source IP.
    * `syn_limit` - The limit of syn of abnormal connection check.
    * `syn_rate` - The percentage of syn in ack of abnormal connection check.
    * `tcp_mbps_limit` - The limit of TCP traffic.
    * `udp_mbps_limit` - The limit of UDP traffic rate.
  * `name` - Name of the DDoS policy.
  * `packet_filters` - Message filter options list.
    * `action` - Action of port to take.
    * `d_end_port` - End port of the destination.
    * `d_start_port` - Start port of the destination.
    * `depth` - The depth of match.
    * `is_include` - Indicate whether to include the key word/regular expression or not.
    * `match_begin` - Indicate whether to check load or not.
    * `match_str` - The key word or regular expression.
    * `match_type` - Match type.
    * `offset` - The offset of match.
    * `pkt_length_max` - The max length of the packet.
    * `pkt_length_min` - The minimum length of the packet.
    * `protocol` - Protocol.
    * `s_end_port` - End port of the source.
    * `s_start_port` - Start port of the source.
  * `policy_id` - Id of policy.
  * `port_limits` - Port limits of abnormal check of the DDoS policy.
    * `action` - Action of port to take.
    * `end_port` - End port.
    * `kind` - The type of forbidden port, and valid values are 0, 1, 2. 0 for destination port, 1 for source port and 2 for both destination and source posts.
    * `protocol` - Protocol.
    * `start_port` - Start port.
  * `scene_id` - Id of scene that the DDoS policy works for.
  * `water_prints` - Water print policy options, and only support one water print policy at most.
    * `auto_remove` - Indicate whether to auto-remove the water print or not.
    * `offset` - The offset of water print.
    * `open_switch` - Indicate whether to open water print or not.
    * `tcp_port_list` - Port range of TCP.
    * `udp_port_list` - Port range of TCP.


