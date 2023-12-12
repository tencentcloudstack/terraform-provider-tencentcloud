Use this data source to query detailed information of tse gateway_nodes

Example Usage

```hcl
data "tencentcloud_tse_gateway_nodes" "gateway_nodes" {
  gateway_id = "gateway-ddbb709b"
  group_id   = "group-013c0d8e"
}
```