Provides a resource to create a tse cngw_network

Example Usage

```hcl
resource "tencentcloud_tse_cngw_network" "cngw_network" {
  gateway_id                 = "gateway-cf8c99c3"
  group_id                   = "group-a160d123"
  internet_address_version   = "IPV4"
  internet_pay_mode          = "BANDWIDTH"
  description                = "des-test1"
  internet_max_bandwidth_out = 1
  master_zone_id             = "ap-guangzhou-3"
  multi_zone_flag            = true
  sla_type                   = "clb.c2.medium"
  slave_zone_id              = "ap-guangzhou-4"
}
```
