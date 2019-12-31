---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dayu_ddos_policy"
sidebar_current: "docs-tencentcloud-resource-dayu_ddos_policy"
description: |-
  Use this resource to create dayu DDoS policy
---

# tencentcloud_dayu_ddos_policy

Use this resource to create dayu DDoS policy

## Example Usage

```hcl
resource "tencentcloud_dayu_ddos_policy" "test_policy" {
  resource_type = "bgpip"
  name          = "tf_test_policy"

  drop_options {
    drop_tcp           = true
    drop_udp           = true
    drop_icmp          = true
    drop_other         = true
    drop_abroad        = true
    check_sync_conn    = true
    source_new_limit   = 100
    dst_new_limit      = 100
    source_conn_limit  = 100
    dst_conn_limit     = 100
    tcp_mbps_limit     = 100
    udp_mbps_limit     = 100
    icmp_mbps_limit    = 100
    other_mbps_limit   = 100
    bad_conn_threshold = 100
    null_conn_enable   = true
    conn_timeout       = 500
    syn_rate           = 50
    syn_limit          = 100
  }

  black_white_ips {
    ip   = "1.1.1.1"
    type = "black"
  }

  port_limits {
    start_port = "2000"
    end_port   = "2500"
    protocol   = "all"
    action     = "drop"
    kind       = 1
  }

  packet_filters {
    protocol       = "tcp"
    action         = "drop"
    d_start_port   = 1000
    d_end_port     = 1500
    s_start_port   = 2000
    s_end_port     = 2500
    pkt_length_max = 1400
    pkt_length_min = 1000
    is_include     = true
    match_begin    = "begin_l5"
    match_type     = "pcre"
    depth          = 1000
    offset         = 500
  }

  water_prints {
    tcp_port_list = ["2000-3000", "3500-4000"]
    udp_port_list = ["5000-6000"]
    offset        = 50
    auto_remove   = true
    open_switch   = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `drop_options` - (Required) Option list of abnormal check of the DDos policy, should set at least one policy.
* `name` - (Required) Name of the DDoS policy. Length should between 1 and 32.
* `resource_type` - (Required, ForceNew) Type of the resource that the DDoS policy works for, valid values are `bgpip`, `bgp`, `bgp-multip`, `net`.
* `black_white_ips` - (Optional) Black and white ip list.
* `packet_filters` - (Optional) Message filter options list.
* `port_limits` - (Optional) Port limits of abnormal check of the DDos policy.
* `water_prints` - (Optional) Water print policy options, and only support one water print policy at most.

The `black_white_ips` object supports the following:

* `ip` - (Optional) Ip.
* `type` - (Optional) Type of the ip, and valid values are `black` and `white`.

The `drop_options` object supports the following:

* `bad_conn_threshold` - (Required) The number of new connections based on destination IP that trigger suppression of connections, and valid value is range from 0 to 4294967295.
* `check_sync_conn` - (Required) Indicate whether to check null connection or not.
* `conn_timeout` - (Required) Connection timeout of abnormal connection check, and valid value is range from 0 to 65535.
* `drop_abroad` - (Required) Indicate whether to drop abroad traffic or not.
* `drop_icmp` - (Required) Indicate whether to drop ICMP protocol or not.
* `drop_other` - (Required) Indicate whether to drop other protocols(exclude TCP/UDP/ICMP) or not.
* `drop_tcp` - (Required) Indicate whether to drop TCP protocol or not.
* `drop_udp` - (Required) Indicate to drop UDP protocol or not.
* `dst_conn_limit` - (Required) The limit of concurrent connections based on destination IP, and valid value is range from 0 to 4294967295.
* `dst_new_limit` - (Required) The limit of new connections based on destination IP, and valid value is range from 0 to 4294967295.
* `icmp_mbps_limit` - (Required) The limit of ICMP traffic rate, and valid value is range from 0 to 4294967295(Mbps).
* `null_conn_enable` - (Required) Indicate to enable null connection or not.
* `other_mbps_limit` - (Required) The limit of other protocols(exclude TCP/UDP/ICMP) traffic rate, and valid value is range from 0 to 4294967295(Mbps).
* `source_conn_limit` - (Required) The limit of concurrent connections based on source IP, and valid value is range from 0 to 4294967295.
* `source_new_limit` - (Required) The limit of new connections based on source IP, and valid value is range from 0 to 4294967295.
* `syn_limit` - (Required) The limit of syn of abnormal connection check, and valid value is range from 0 to 100.
* `tcp_mbps_limit` - (Required) The limit of TCP traffic, and valid value is range from 0 to 4294967295(Mbps).
* `udp_mbps_limit` - (Required) The limit of UDP traffic rate, and valid value is range from 0 to 4294967295(Mbps).
* `syn_rate` - (Optional) The percentage of syn in ack of abnormal connection check, and valid value is range from 0 to 100.

The `packet_filters` object supports the following:

* `action` - (Optional) Action of port to take, valid values area `drop`(drop the packet), `drop_black`(drop the packet and black the ip),`drop_rst`(drop the packet and disconnect),`drop_black_rst`(drop the packet, black the ip and disconnect),`transmit`(transmit the packet).
* `d_end_port` - (Optional) End port of the destination, valid value is range from 0 to 65535. It must be greater than `d_start_port`.
* `d_start_port` - (Optional) Start port of the destination, valid value is range from 0 to 65535.
* `depth` - (Optional) The depth of match, and valid value is range from 0 to 1500.
* `is_include` - (Optional) Indicate whether to include the key word/regular expression or not.
* `match_begin` - (Optional) Indicate whether to check load or not, `begin_l5` means to match and `no_match` means not.
* `match_str` - (Optional) The key word or regular expression.
* `match_type` - (Optional) Match type, valid values are `sunday` and `pcre`, `sunday` means key word match while `pcre` means regular match.
* `offset` - (Optional) The offset of match, and valid value is range from 0 to 1500.
* `pkt_length_max` - (Optional) The max length of the packet, and valid value is range from 0 to 1500(Mbps). It must be greater than `pkt_length_min`.
* `pkt_length_min` - (Optional) The minimum length of the packet, and valid value is range from 0 to 1500(Mbps).
* `protocol` - (Optional) Protocol, valid values are `tcp`, `udp`, `icmp`, `all`.
* `s_end_port` - (Optional) End port of the source, valid value is range from 0 to 65535. It must be greater than `s_start_port`.
* `s_start_port` - (Optional) Start port of the source, valid value is range from 0 to 65535.

The `port_limits` object supports the following:

* `action` - (Optional) Action of port to take, valid values area `drop`, `transmit`.
* `end_port` - (Optional) End port, valid value is range from 0 to 65535. It must be greater than `d_start_port`.
* `kind` - (Optional) The type of forbidden port, and valid values are 0, 1, 2. 0 for destination ports make effect, 1 for source ports make effect. 2 for both destination and source ports.
* `protocol` - (Optional) Protocol, valid values are `tcp`, `udp`, `icmp`, `all`.
* `start_port` - (Optional) Start port, valid value is range from 0 to 65535.

The `water_prints` object supports the following:

* `auto_remove` - (Optional) Indicate whether to auto-remove the water print or not.
* `offset` - (Optional) The offset of water print, and valid value is range from 0 to 100.
* `open_switch` - (Optional) Indicate whether to open water print or not.
* `tcp_port_list` - (Optional) Port range of TCP, the format is like `2000-3000`.
* `udp_port_list` - (Optional) Port range of TCP, the format is like `2000-3000`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - Create time of the DDos policy.
* `policy_id` - Id of policy.
* `scene_id` - Id of scene that the DDos policy works for.


