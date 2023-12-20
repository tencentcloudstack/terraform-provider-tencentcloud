Provides a resource to create a vpc peer_connect_reject_operation

Example Usage

```hcl
resource "tencentcloud_vpc_peer_connect_reject_operation" "peer_connect_reject_operation" {
  peering_connection_id = "pcx-abced"
}
```

Import

vpc peer_connect_reject_operation can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_peer_connect_reject_operation.peer_connect_reject_operation peer_connect_reject_operation_id
```
