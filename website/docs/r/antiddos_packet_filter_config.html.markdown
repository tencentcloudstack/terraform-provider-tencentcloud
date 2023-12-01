---
subcategory: "Anti-DDoS(DayuV2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_antiddos_packet_filter_config"
sidebar_current: "docs-tencentcloud-resource-antiddos_packet_filter_config"
description: |-
  Provides a resource to create a antiddos packet filter config
---

# tencentcloud_antiddos_packet_filter_config

Provides a resource to create a antiddos packet filter config

## Example Usage

```hcl
resource "tencentcloud_antiddos_packet_filter_config" "packet_filter_config" {
  instance_id = "bgp-00000ry7"
  packet_filter_config {
    action      = "drop"
    depth       = 1
    dport_start = 80
    dport_end   = 80
    is_not      = 0
    match_begin = "begin_l5"
    match_type  = "pcre"
    offset      = 1
    pktlen_min  = 1400
    pktlen_max  = 1400
    protocol    = "all"
    sport_start = 8080
    sport_end   = 8080
    str         = "a"
  }
} s
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) resource id.
* `packet_filter_config` - (Required, List, ForceNew) Feature filtering configuration.

The `packet_filter_config` object supports the following:

* `action` - (Required, String) Action, value [drop (discard) transmit (release) drop_black (discard and pull black) drop_rst (intercept) drop_black_rst (intercept and pull black) forward (continue protection)].
* `dport_end` - (Required, Int) end destination port, ranging from 0 to 65535.
* `dport_start` - (Required, Int) Starting destination port, ranging from 0 to 65535.
* `pktlen_max` - (Required, Int) The maximum message length, ranging from 1 to 1500, must be greater than or equal to the minimum message length.
* `pktlen_min` - (Required, Int) Minimum message length, ranging from 1 to 1500.
* `protocol` - (Required, String) Protocol, value [TCP udp icmp all].
* `sport_end` - (Required, Int) End source port, values range from 1 to 65535, must be greater than or equal to the start source port.
* `sport_start` - (Required, Int) Starting source port, ranging from 0 to 65535.
* `depth2` - (Optional, Int) The second detection depth starting from the second detection position, with a value of [01500].
* `depth` - (Optional, Int) The detection depth starting from the detection position, with a value of [0-1500].
* `is_not2` - (Optional, Int) Whether the second detection includes detection values, with a value of [0 (inclusive) and 1 (exclusive)].
* `is_not` - (Optional, Int) Whether to include detection values, with a value of [0 (inclusive) and 1 (exclusive)].
* `match_begin2` - (Optional, String) Second detection position, value [begin_l5 (load) no_match (mismatch)].
* `match_begin` - (Optional, String) Detection position, value [begin_l3 (IP header) begin_l4 (TCP/UDP header) begin_l5 (T payload) no_match (mismatch)].
* `match_logic` - (Optional, String) When there is a second detection condition, the AND or relationship with the first detection condition, with the value [and (and relationship) none (fill in this value when there is no second detection condition)].
* `match_type2` - (Optional, String) The second detection type, with a value of [Sunday (keyword) pcre (regular expression)].
* `match_type` - (Optional, String) Detection type, value [Sunday (keyword) pcre (regular expression)].
* `offset2` - (Optional, Int) The offset from the second detection position, with a value range of [0, Depth2].
* `offset` - (Optional, Int) The offset from the detection position, with a value range of [0, Depth].
* `pkt_len_gt` - (Optional, Int) Greater than message length, value 1+.
* `str2` - (Optional, String) key string or regular expression, value [When the detection type is Sunday, please fill in the string or hexadecimal bytecode, for example, x313233 corresponds to the hexadecimal word&gt;section code of the string &#39;123&#39;; when the detection type is pcre, please fill in the regular expression character string;].
* `str` - (Optional, String) Detection value, key string or regular expression, value [When the detection type is Sunday, please fill in the string or hexadecimal bytecode, for example, x313233 corresponds to the hexadecimal word&gt;section code of the string &#39;123&#39;; when the detection type is pcre, please fill in the regular expression character string;].

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

antiddos packet_filter_config can be imported using the id, e.g.

```
terraform import tencentcloud_antiddos_packet_filter_config.packet_filter_config packet_filter_config_id
```

