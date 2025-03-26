Provides a resource to create a DC extra config

Example Usage

```hcl
resource "tencentcloud_dcx_extra_config" "example" {
  direct_connect_tunnel_id = "dcx-4z49tnws"
  vlan                     = 123
  tencent_address          = "10.3.191.73/29"
  tencent_backup_address   = "10.3.191.72/29"
  customer_address         = "10.3.191.74/29"
  bandwidth                = 100
  enable_bgp_community     = false
  bfd_enable               = 1
  nqa_enable               = 0
  bgp_peer {
    asn      = 65101
    auth_key = "test123"
  }
  bfd_info {
    probe_failed_times = 3
    interval           = 2000
  }
  nqa_info {
    probe_failed_times = -1
    interval           = -1
    destination_ip     = "0.0.0.0"
  }
  ipv6_enable  = 0
  jumbo_enable = 0
}
```

Import

DC extra config can be imported using the id, e.g.

```
terraform import tencentcloud_dcx_extra_config.example dcx-4z49tnws
```