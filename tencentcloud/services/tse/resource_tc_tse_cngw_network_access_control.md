Provides a resource to create a tse cngw_network_access_control

Example Usage

```hcl
resource "tencentcloud_tse_cngw_network_access_control" "cngw_network_access_control" {
  gateway_id                 = "gateway-cf8c99c3"
  group_id                   = "group-a160d123"
  network_id                 = "network-372b1e84"
  access_control {
    mode            = "Whitelist"
    cidr_white_list = ["1.1.1.0"]
  }
}
```

Import

tse cngw_route_rate_limit can be imported using the id, e.g.

```
terraform import tencentcloud_tse_cngw_network_access_control.cngw_network_access_control gatewayId#groupId#networkId
```