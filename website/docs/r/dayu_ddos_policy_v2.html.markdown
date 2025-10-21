---
subcategory: "Anti-DDoS(DayuV2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dayu_ddos_policy_v2"
sidebar_current: "docs-tencentcloud-resource-dayu_ddos_policy_v2"
description: |-
  Use this resource to create dayu DDoS policy v2
---

# tencentcloud_dayu_ddos_policy_v2

Use this resource to create dayu DDoS policy v2

## Example Usage

```hcl
resource "tencentcloud_dayu_ddos_policy_v2" "ddos_v2" {
  resource_id    = "bgpip-000004xf"
  business       = "bgpip"
  ddos_threshold = "100"
  ddos_level     = "low"
  black_white_ips {
    ip      = "1.2.3.4"
    ip_type = "black"
  }
  acls {
    action           = "transmit"
    d_port_start     = 1
    d_port_end       = 10
    s_port_start     = 10
    s_port_end       = 20
    priority         = 9
    forward_protocol = "all"
  }
  protocol_block_config {
    drop_icmp  = 1
    drop_tcp   = 0
    drop_udp   = 0
    drop_other = 0
  }
  ddos_connect_limit {
    sd_new_limit       = 10
    sd_conn_limit      = 11
    dst_new_limit      = 20
    dst_conn_limit     = 21
    bad_conn_threshold = 30
    syn_rate           = 10
    syn_limit          = 20
    conn_timeout       = 30
    null_conn_enable   = 1
  }
  ddos_ai = "on"
  ddos_geo_ip_block_config {
    action      = "drop"
    area_list   = ["100001"]
    region_type = "customized"
  }
  ddos_speed_limit_config {
    protocol_list = "TCP"
    dst_port_list = "10"
    mode          = 1
    packet_rate   = 10
    bandwidth     = 20
  }
  packet_filters {
    action       = "drop"
    protocol     = "all"
    s_port_start = 10
    s_port_end   = 10
    d_port_start = 20
    d_port_end   = 20
    pktlen_min   = 30
    pktlen_max   = 30
    str          = "12"
    str2         = "30"
    match_logic  = "and"
    match_type   = "pcre"
    match_type2  = "pcre"
    match_begin  = "begin_l3"
    match_begin2 = "begin_l3"
    depth        = 2
    depth2       = 3
    offset       = 1
    offset2      = 2
    is_not       = 0
    is_not2      = 0
  }
  water_print_config {
    offset      = 1
    open_status = 1
    listeners {
      frontend_port     = 90
      forward_protocol  = "TCP"
      frontend_port_end = 90
    }
    verify = "checkall"
  }
}
```

## Argument Reference

The following arguments are supported:

* `resource_id` - (Required, String, ForceNew) The ID of the resource instance.
* `acls` - (Optional, List) Port ACL policy for DDoS protection.
* `black_white_ips` - (Optional, List) DDoS-protected IP blacklist and whitelist.
* `business` - (Optional, String) Business of resource instance. bgpip indicates anti-anti-ip ip; bgp means exclusive package; bgp-multip means shared packet; net indicates anti-anti-ip pro version.
* `ddos_ai` - (Optional, String) AI protection switch, take the value [`on`, `off`].
* `ddos_connect_limit` - (Optional, List) DDoS connection suppression options.
* `ddos_geo_ip_block_config` - (Optional, List) DDoS-protected area block configuration.
* `ddos_level` - (Optional, String) Protection class, value [`low`, `middle`, `high`].
* `ddos_speed_limit_config` - (Optional, List) Access speed limit configuration for DDoS protection.
* `ddos_threshold` - (Optional, Int) DDoS cleaning threshold, value[0, 60, 80, 100, 150, 200, 250, 300, 400, 500, 700, 1000]; When the value is set to 0, it means that the default value is adopted.
* `packet_filters` - (Optional, List) Feature filtering rules for DDoS protection.
* `protocol_block_config` - (Optional, List) Protocol block configuration for DDoS protection.
* `water_print_config` - (Optional, List) Water print config.

The `acls` object supports the following:

* `action` - (Required, String) Action, optional values: drop, transmit, forward.
* `d_port_end` - (Required, Int) The destination port ends, and the value range is 0~65535.
* `d_port_start` - (Required, Int) The destination port starts, and the value range is 0~65535.
* `forward_protocol` - (Required, String) Protocol type, desirable values tcp, udp, all.
* `priority` - (Required, Int) Policy priority, the lower the number, the higher the level, the higher the rule matches, taking a value of 1-1000.Note: This field may return null, indicating that a valid value could not be retrieved.
* `s_port_end` - (Required, Int) The source port ends, and the acceptable value ranges from 0 to 65535.
* `s_port_start` - (Required, Int) The source port starts, and the value range is 0~65535.

The `black_white_ips` object supports the following:

* `ip_type` - (Required, String) IP type, value [`black`(blacklist IP), `white` (whitelist IP)].
* `ip` - (Required, String) Ip of resource instance.

The `ddos_connect_limit` object supports the following:

* `bad_conn_threshold` - (Required, Int) Based on connection suppression trigger threshold, value range [0,4294967295].
* `conn_timeout` - (Required, Int) Abnormal connection detection condition, connection timeout, value range [0,65535].
* `dst_conn_limit` - (Required, Int) Concurrent connection control based on destination IP+ destination port.
* `dst_new_limit` - (Required, Int) Limit on the number of news per second based on the destination IP.
* `null_conn_enable` - (Required, Int) Abnormal connection detection conditions, empty connection guard switch, value range[0,1].
* `sd_conn_limit` - (Required, Int) Concurrent connection control based on source IP + destination IP.
* `sd_new_limit` - (Required, Int) The limit on the number of news per second based on source IP + destination IP.
* `syn_limit` - (Required, Int) Anomaly connection detection condition, syn threshold, value range [0,100].
* `syn_rate` - (Required, Int) Anomalous connection detection condition, percentage of syn ack, value range [0,100].

The `ddos_geo_ip_block_config` object supports the following:

* `action` - (Required, String) Block action, take the value [`drop`, `trans`].
* `area_list` - (Required, List) When the RegionType is customized, the AreaList must be filled in, and a maximum of 128 must be filled in.
* `region_type` - (Required, String) Zone type, value [oversea (overseas),china (domestic),customized (custom region)].

The `ddos_speed_limit_config` object supports the following:

* `bandwidth` - (Required, Int) Bandwidth bps.
* `dst_port_list` - (Required, String) List of port ranges, up to 8, multiple; Separated, the range is represented with -; this port range must be filled in; fill in the style 1:0-65535, style 2:80; 443; 1000-2000.
* `mode` - (Required, Int) Speed limit mode, take the value [1 (speed limit based on source IP),2 (speed limit based on destination port)].
* `packet_rate` - (Required, Int) Packet rate pps.
* `protocol_list` - (Required, String) IP protocol numbers, take the value[ ALL (all protocols),TCP (tcp protocol),UDP (udp protocol),SMP (smp protocol),1; 2-100 (custom protocol number range, up to 8)].

The `listeners` object of `water_print_config` supports the following:

* `forward_protocol` - (Required, String) Forwarding protocol, value [TCP, UDP].
* `frontend_port_end` - (Required, Int) Upper limit of forwarding listening port. Values: [1-65535].
* `frontend_port` - (Required, Int) Lower limit of forwarding listening port. Values: [1-65535].

The `packet_filters` object supports the following:

* `action` - (Required, String) Action, take the value [drop,transmit,drop_black (discard and black out),drop_rst (Interception),drop_black_rst (intercept and block),forward].
* `d_port_end` - (Required, Int) The end destination port, take the value 1~65535, which must be greater than or equal to the starting destination port.
* `d_port_start` - (Required, Int) From the destination port, take the value 0~65535.
* `depth2` - (Required, Int) Second detection depth starting from the second detection position, value [0,1500].
* `depth` - (Required, Int) Detection depth from the detection position, value [0,1500].
* `is_not2` - (Required, Int) Whether the second detection contains the detected value, the value [0 (included),1 (not included)].
* `is_not` - (Required, Int) Whether to include the detected value, take the value [0 (included),1 (not included)].
* `match_begin2` - (Required, String) The second detection position. take the value [begin_l3 (IP header),begin_l4 (TCP/UDP header),begin_l5 (T load), no_match (mismatch)].
* `match_begin` - (Required, String) Detect position, take the value [begin_l3 (IP header),begin_l4 (TCP/UDP header),begin_l5 (T load), no_match (mismatch)].
* `match_logic` - (Required, String) When there is a second detection condition, the and/or relationship with the first detection condition, takes the value [And (and relationship),none (fill in this value when there is no second detection condition)].
* `match_type2` - (Required, String) The second type of detection, takes the value [sunday (keyword),pcre (regular expression)].
* `match_type` - (Required, String) Detection type, value [sunday (keyword),pcre (regular expression)].
* `offset2` - (Required, Int) Offset from the second detection position, value range [0,Depth2].
* `offset` - (Required, Int) Offset from detection position, value range [0, Depth].
* `pktlen_max` - (Required, Int) The maximum message length, taken from 1 to 1500, must be greater than or equal to the minimum message length.
* `pktlen_min` - (Required, Int) Minimum message length, 1-1500.
* `protocol` - (Required, String) Protocol, value [tcp udp icmp all].
* `s_port_end` - (Required, Int) End source port, take the value 1~65535, must be greater than or equal to the starting source port.
* `s_port_start` - (Required, Int) Start the source port, take the value 0~65535.
* `str2` - (Required, String) The second detection value, the key string or regular expression, takes the value [When the detection type is sunday, please fill in the string or hexadecimal bytecode, for example 13233 corresponds to the hexadecimal bytecode of the string `123`;When the detection type is pcre, please fill in the regular expression string;].
* `str` - (Required, String) Detect values, key strings or regular expressions, take the value [When the detection type is sunday, please fill in the string or hexadecimal bytecode, for example 13233 corresponds to the hexadecimal bytecode of the string `123`;When the detection type is pcre, please fill in the regular expression string;].

The `protocol_block_config` object supports the following:

* `drop_icmp` - (Required, Int) ICMP block, value [0 (block off), 1 (block on)].
* `drop_other` - (Required, Int) Other block, value [0 (block off), 1 (block on)].
* `drop_tcp` - (Required, Int) TCP block, value [0 (block off), 1 (block on)].
* `drop_udp` - (Required, Int) UDP block, value [0 (block off), 1 (block on)].

The `water_print_config` object supports the following:

* `listeners` - (Required, List) List of forwarding listeners to which the watermark belongs.
* `offset` - (Required, Int) Watermark offset, value range: [0-100].
* `open_status` - (Required, Int) Whether it is enabled, value [0 (manual open), 1 (immediate operation)].
* `verify` - (Required, String) Watermark check mode, value [`checkall`(normal mode), `shortfpcheckall`(simplified mode)].

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



