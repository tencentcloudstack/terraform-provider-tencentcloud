Provides a resource to create a vpc peer_connect_manager

Example Usage

```hcl
resource "tencentcloud_vpc_peer_connect_manager" "peer_connect_manager" {
  source_vpc_id = "vpc-abcdef"
  peering_connection_name = "name"
  destination_vpc_id = "vpc-abc1234"
  destination_uin = "12345678"
  destination_region = "ap-beijing"
  bandwidth = 100
  type = "VPC_PEER"
  charge_type = "POSTPAID_BY_DAY_MAX"
  qos_level = "AU"
}
```

Import

vpc peer_connect_manager can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_peer_connect_manager.peer_connect_manager peer_connect_manager_id
```
