---
subcategory: "Teo"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_ddos_policy"
sidebar_current: "docs-tencentcloud-resource-teo_ddos_policy"
description: |-
  Provides a resource to create a teo ddosPolicy
---

# tencentcloud_teo_ddos_policy

Provides a resource to create a teo ddosPolicy

## Example Usage

```hcl
resource "tencentcloud_teo_ddos_policy" "ddosPolicy" {
  zone_id   = ""
  policy_id = ""
  ddos_rule {
    switch         = ""
    udp_shard_open = ""
    ddos_status_info {
      ply_level = ""
    }
    ddos_geo_ip {
      region_id = ""
      switch    = ""
    }
    ddos_allow_block {
      switch = ""
      user_allow_block_ip {
        ip    = ""
        mask  = ""
        type  = ""
        ip2   = ""
        mask2 = ""
      }
    }
    ddos_anti_ply {
      drop_tcp                  = ""
      drop_udp                  = ""
      drop_icmp                 = ""
      drop_other                = ""
      source_create_limit       = ""
      source_connect_limit      = ""
      destination_create_limit  = ""
      destination_connect_limit = ""
      abnormal_connect_num      = ""
      abnormal_syn_ratio        = ""
      abnormal_syn_num          = ""
      connect_timeout           = ""
      empty_connect_protect     = ""
      udp_shard                 = ""
    }
    ddos_packet_filter {
      switch = ""
      packet_filter {
        action       = ""
        protocol     = ""
        dport_start  = ""
        dport_end    = ""
        packet_min   = ""
        packet_max   = ""
        sport_start  = ""
        sport_end    = ""
        match_type   = ""
        is_not       = ""
        offset       = ""
        depth        = ""
        match_begin  = ""
        str          = ""
        match_type2  = ""
        is_not2      = ""
        offset2      = ""
        depth2       = ""
        match_begin2 = ""
        str2         = ""
        match_logic  = ""
      }
    }
    ddos_acl {
      switch = ""
      acl {
        dport_end   = ""
        dport_start = ""
        sport_end   = ""
        sport_start = ""
        protocol    = ""
        action      = ""
        default     = ""
      }
    }

  }
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, Int) Policy ID.
* `zone_id` - (Required, String) Site ID.
* `ddos_rule` - (Optional, List) DDoS Configuration of the zone.

The `acl` object supports the following:

* `action` - (Optional, String) Action to take. Valid values: `drop`, `transmit`, `forward`.
* `default` - (Optional, Int) Whether it is default configuration. Valid value:- 0: custom configuration.- 1: default configuration.
* `dport_end` - (Optional, Int) End of the dest port range. Valid value range: 0-65535.
* `dport_start` - (Optional, Int) Start of the dest port range. Valid value range: 0-65535.
* `protocol` - (Optional, String) Valid values: `tcp`, `udp`, `all`.
* `sport_end` - (Optional, Int) End of the source port range. Valid value range: 0-65535.
* `sport_start` - (Optional, Int) Start of the source port range. Valid value range: 0-65535.

The `ddos_acl` object supports the following:

* `acl` - (Optional, List) DDoS ACL rule configuration detail.
* `switch` - (Optional, String) - on: Enable. `Acl` parameter is require.- off: Disable.

The `ddos_allow_block` object supports the following:

* `switch` - (Optional, String) - on: Enable. `UserAllowBlockIp` parameter is required.- off: Disable.
* `user_allow_block_ip` - (Optional, List) DDoS black-white list detail.

The `ddos_anti_ply` object supports the following:

* `abnormal_connect_num` - (Required, Int) Abnormal connections threshold. Valid value range: 0-4294967295.
* `abnormal_syn_num` - (Required, Int) Abnormal syn packet number threshold. Valid value range: 0-65535.
* `abnormal_syn_ratio` - (Required, Int) Abnormal syn packet ratio threshold. Valid value range: 0-100.
* `connect_timeout` - (Required, Int) Connection timeout detection per second. Valid value range: 0-65535.
* `destination_connect_limit` - (Required, Int) Limitation of connections to dest port. Valid value range: 0-4294967295.
* `destination_create_limit` - (Required, Int) Limitation of new connection to dest port per second. Valid value range: 0-4294967295.
* `drop_icmp` - (Required, String) Block ICMP protocol. Valid values: `on`, `off`.
* `drop_other` - (Required, String) Block other protocols. Valid values: `on`, `off`.
* `drop_tcp` - (Required, String) Block TCP protocol. Valid values: `on`, `off`.
* `drop_udp` - (Required, String) Block UDP protocol. Valid values: `on`, `off`.
* `empty_connect_protect` - (Required, String) Empty connection protection switch. Valid values: `on`, `off`.
* `source_connect_limit` - (Required, Int) Limitation of connections to origin website. Valid value range: 0-4294967295.
* `source_create_limit` - (Required, Int) Limitation of new connection to origin website per second. Valid value range: 0-4294967295.
* `udp_shard` - (Optional, String) UDP shard protection switch. Valid values: `on`, `off`.

The `ddos_geo_ip` object supports the following:

* `region_id` - (Optional, Set) Region ID. See details in data source `security_policy_regions`.
* `switch` - (Optional, String) - on: Enable.- off: Disable.

The `ddos_packet_filter` object supports the following:

* `packet_filter` - (Optional, List) DDoS feature filtering configuration detail.
* `switch` - (Optional, String) - on: Enable. `ModifyDDoSPolicy` parameter is required.- off: Disable.

The `ddos_rule` object supports the following:

* `ddos_acl` - (Optional, List) DDoS ACL rule configuration.
* `ddos_allow_block` - (Optional, List) DDoS black-white list.
* `ddos_anti_ply` - (Optional, List) DDoS protocol and connection protection.
* `ddos_geo_ip` - (Optional, List) DDoS Protection by Geo Info.
* `ddos_packet_filter` - (Optional, List) DDoS feature filtering configuration.
* `ddos_status_info` - (Optional, List) DDoS protection level.
* `switch` - (Optional, String) DDoS protection switch. Valid values:- on: Enable.- off: Disable.
* `udp_shard_open` - (Optional, String) UDP shard switch. Valid values:- on: Enable.- off: Disable.

The `ddos_status_info` object supports the following:

* `ply_level` - (Required, String) Policy level. Valid values:- low: loose.- middle: moderate.- high: strict.

The `packet_filter` object supports the following:

* `action` - (Optional, String) Action to take. Valid values: `drop`, `transmit`, `drop_block`, `forward`.
* `depth2` - (Optional, Int) Packet character depth to check of feature 2. Valid value range: 1-1500.
* `depth` - (Optional, Int) Packet character depth to check of feature 1. Valid value range: 1-1500.
* `dport_end` - (Optional, Int) End of the dest port range. Valid value range: 0-65535.
* `dport_start` - (Optional, Int) Start of the dest port range. Valid value range: 0-65535.
* `is_not2` - (Optional, Int) Negate the match condition of feature 2. Valid values:- 0: match.- 1: not match.
* `is_not` - (Optional, Int) Negate the match condition of feature 1. Valid values:- 0: match.- 1: not match.
* `match_begin2` - (Optional, String) Packet layer for matching begin of feature 2. Valid values:- begin_l5: matching from packet payload.- begin_l4: matching from TCP/UDP header.- begin_l3: matching from IP header.
* `match_begin` - (Optional, String) Packet layer for matching begin of feature 1. Valid values:- begin_l5: matching from packet payload.- begin_l4: matching from TCP/UDP header.- begin_l3: matching from IP header.
* `match_logic` - (Optional, String) Relation between multi features. Valid values: `and`, `or`, `none` (only feature 1 is used).
* `match_type2` - (Optional, String) Match type of feature 2. Valid values:- pcre: regex expression.- sunday: string match.
* `match_type` - (Optional, String) Match type of feature 1. Valid values:- pcre: regex expression.- sunday: string match.
* `offset2` - (Optional, Int) Offset of feature 2. Valid value range: 1-1500.
* `offset` - (Optional, Int) Offset of feature 1. Valid value range: 1-1500.
* `packet_max` - (Optional, Int) Max packet size. Valid value range: 0-1500.
* `packet_min` - (Optional, Int) Min packet size. Valid value range: 0-1500.
* `protocol` - (Optional, String) Valid value: `tcp`, `udp`, `icmp`, `all`.
* `sport_end` - (Optional, Int) End of the source port range. Valid value range: 0-65535.
* `sport_start` - (Optional, Int) Start of the source port range. Valid value range: 0-65535.
* `str2` - (Optional, String) Regex expression or string to match.
* `str` - (Optional, String) Regex expression or string to match.

The `user_allow_block_ip` object supports the following:

* `type` - (Required, String) Valid values: `block`, `allow`.
* `ip2` - (Optional, String) End of the IP address when setting an IP range.
* `ip` - (Optional, String) Client IP.
* `mask2` - (Optional, Int) IP mask of the end IP address.
* `mask` - (Optional, Int) IP Mask.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

teo ddosPolicy can be imported using the id, e.g.
```
$ terraform import tencentcloud_teo_ddos_policy.ddosPolicy ddosPolicy_id
```

