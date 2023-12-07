Provides a resource to create a antiddos packet filter config

Example Usage

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
}s
```

Import

antiddos packet_filter_config can be imported using the id, e.g.

```
terraform import tencentcloud_antiddos_packet_filter_config.packet_filter_config packet_filter_config_id
```