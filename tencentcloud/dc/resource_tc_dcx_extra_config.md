Provides a resource to create a dc dcx_extra_config

Example Usage

```hcl
resource "tencentcloud_dcx_extra_config" "dcx_extra_config" {
  direct_connect_tunnel_id = "dcx-4z49tnws"
  vlan                     = 123
  bgp_peer {
    asn      = 65101
    auth_key = "test123"

  }
  route_filter_prefixes {
    cidr = "192.168.0.0/24"
  }
  tencent_address        = "192.168.1.1"
  tencent_backup_address = "192.168.1.2"
  customer_address       = "192.168.1.4"
  bandwidth              = 10
  enable_bgp_community   = false
  bfd_enable             = 0
  nqa_enable             = 1
  bfd_info {
    probe_failed_times = 3
    interval           = 100

  }
  nqa_info {
    probe_failed_times = 3
    interval           = 100
    destination_ip     = "192.168.2.2"

  }
  ipv6_enable = 0
  jumbo_enable = 0
}
```

Import

dc dcx_extra_config can be imported using the id, e.g.

```
terraform import tencentcloud_dcx_extra_config.dcx_extra_config dcx_id
```