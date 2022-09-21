---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_ddos_policy"
sidebar_current: "docs-tencentcloud-resource-teo_ddos_policy"
description: |-
  Provides a resource to create a teo ddos_policy
---

# tencentcloud_teo_ddos_policy

Provides a resource to create a teo ddos_policy

## Example Usage

```hcl
resource "tencentcloud_teo_ddos_policy" "ddos_policy" {
  policy_id = 1278
  zone_id   = "zone-2983wizgxqvm"

  ddos_rule {
    switch = "on"

    acl {
      switch = "on"
    }

    allow_block {
      switch = "on"
    }

    anti_ply {
      abnormal_connect_num      = 0
      abnormal_syn_num          = 0
      abnormal_syn_ratio        = 0
      connect_timeout           = 0
      destination_connect_limit = 0
      destination_create_limit  = 0
      drop_icmp                 = "off"
      drop_other                = "off"
      drop_tcp                  = "off"
      drop_udp                  = "off"
      empty_connect_protect     = "off"
      source_connect_limit      = 0
      source_create_limit       = 0
      udp_shard                 = "off"
    }

    geo_ip {
      region_ids = []
      switch     = "on"
    }

    packet_filter {
      switch = "on"
    }

    speed_limit {
      flux_limit    = "0 bps"
      package_limit = "0 pps"
    }

    status_info {
      ply_level = "middle"
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

* `acls` - (Optional, List) DDoS ACL rule configuration detail.
* `switch` - (Optional, String) - `on`: Enable. `Acl` parameter is require.- `off`: Disable.

The `acls` object supports the following:

* `action` - (Optional, String) Action to take. Valid values: `drop`, `transmit`, `forward`.
* `dport_end` - (Optional, Int) End of the dest port range. Valid value range: 0-65535.
* `dport_start` - (Optional, Int) Start of the dest port range. Valid value range: 0-65535.
* `protocol` - (Optional, String) Valid values: `tcp`, `udp`, `all`.
* `sport_end` - (Optional, Int) End of the source port range. Valid value range: 0-65535.
* `sport_start` - (Optional, Int) Start of the source port range. Valid value range: 0-65535.

The `allow_block_ips` object supports the following:

* `type` - (Required, String) Valid values: `block`, `allow`.
* `ip` - (Optional, String) Valid value format:- ip, for example 1.1.1.1- ip range, for example 1.1.1.2-1.1.1.3- network segment, for example 1.2.1.0/24- network segment range, for example 1.2.1.0/24-1.2.2.0/24.

The `allow_block` object supports the following:

* `allow_block_ips` - (Optional, List) DDoS black-white list detail.
* `switch` - (Optional, String) - `on`: Enable. `AllowBlockIps` parameter is required.- `off`: Disable.

The `anti_ply` object supports the following:

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
* `source_connect_limit` - (Required, Int) Limitation of connections to origin site. Valid value range: 0-4294967295.
* `source_create_limit` - (Required, Int) Limitation of new connection to origin site per second. Valid value range: 0-4294967295.
* `udp_shard` - (Optional, String) UDP shard protection switch. Valid values: `on`, `off`.

The `ddos_rule` object supports the following:

* `acl` - (Optional, List) DDoS ACL rule configuration.
* `allow_block` - (Optional, List) DDoS black-white list.
* `anti_ply` - (Optional, List) DDoS protocol and connection protection.
* `geo_ip` - (Optional, List) DDoS Protection by Geo Info.
* `packet_filter` - (Optional, List) DDoS feature filtering configuration.
* `speed_limit` - (Optional, List) DDoS access origin site speed limit configuration.
* `status_info` - (Optional, List) DDoS protection level.
* `switch` - (Optional, String) DDoS protection switch. Valid values:- `on`: Enable.- `off`: Disable.

The `geo_ip` object supports the following:

* `region_ids` - (Optional, Set) Region ID. See details in data source `security_policy_regions`.
* `switch` - (Optional, String) - `on`: Enable.- `off`: Disable.

The `packet_filter` object supports the following:

* `packet_filters` - (Optional, List) DDoS feature filtering configuration detail.
* `switch` - (Optional, String) - `on`: Enable. `PacketFilters` parameter is required.- `off`: Disable.

The `packet_filters` object supports the following:

* `action` - (Optional, String) Action to take. Valid values: `drop`, `transmit`, `drop_block`, `forward`.
* `depth2` - (Optional, Int) Packet character depth to check of feature 2. Valid value range: 1-1500.
* `depth` - (Optional, Int) Packet character depth to check of feature 1. Valid value range: 1-1500.
* `dport_end` - (Optional, Int) End of the dest port range. Valid value range: 0-65535.
* `dport_start` - (Optional, Int) Start of the dest port range. Valid value range: 0-65535.
* `is_not2` - (Optional, Int) Negate the match condition of feature 2. Valid values:- `0`: match.- `1`: not match.
* `is_not` - (Optional, Int) Negate the match condition of feature 1. Valid values:- `0`: match.- `1`: not match.
* `match_begin2` - (Optional, String) Packet layer for matching begin of feature 2. Valid values:- `begin_l5`: matching from packet payload.- `begin_l4`: matching from TCP/UDP header.- `begin_l3`: matching from IP header.
* `match_begin` - (Optional, String) Packet layer for matching begin of feature 1. Valid values:- `begin_l5`: matching from packet payload.- `begin_l4`: matching from TCP/UDP header.- `begin_l3`: matching from IP header.
* `match_logic` - (Optional, String) Relation between multi features. Valid values: `and`, `or`, `none` (only feature 1 is used).
* `match_type2` - (Optional, String) Match type of feature 2. Valid values:- `pcre`: regex expression.- `sunday`: string match.
* `match_type` - (Optional, String) Match type of feature 1. Valid values:- `pcre`: regex expression.- `sunday`: string match.
* `offset2` - (Optional, Int) Offset of feature 2. Valid value range: 1-1500.
* `offset` - (Optional, Int) Offset of feature 1. Valid value range: 1-1500.
* `packet_max` - (Optional, Int) Max packet size. Valid value range: 0-1500.
* `packet_min` - (Optional, Int) Min packet size. Valid value range: 0-1500.
* `protocol` - (Optional, String) Valid value: `tcp`, `udp`, `icmp`, `all`.
* `sport_end` - (Optional, Int) End of the source port range. Valid value range: 0-65535.
* `sport_start` - (Optional, Int) Start of the source port range. Valid value range: 0-65535.
* `str2` - (Optional, String) Regex expression or string to match.
* `str` - (Optional, String) Regex expression or string to match.

The `speed_limit` object supports the following:

* `flux_limit` - (Optional, String) Limit the number of fluxes. Valid range: 1 bps-10000 Gbps, 0 means no limitation, supported units: `pps`,`Kpps`,`Mpps`,`Gpps`.
* `package_limit` - (Optional, String) Limit the number of packages. Valid range: 1 pps-10000 Gpps, 0 means no limitation, supported units: `pps`,`Kpps`,`Mpps`,`Gpps`.

The `status_info` object supports the following:

* `ply_level` - (Required, String) Policy level. Valid values:- `low`: loose.- `middle`: moderate.- `high`: strict.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

teo ddos_policy can be imported using the id, e.g.
```
$ terraform import tencentcloud_teo_ddos_policy.ddos_policy ddosPolicy_id
```

