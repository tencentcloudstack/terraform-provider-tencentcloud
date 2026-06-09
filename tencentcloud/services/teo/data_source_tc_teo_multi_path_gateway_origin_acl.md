Use this data source to query detailed information of TEO multi-path gateway origin acl

Example Usage

Query multi-path gateway origin acl by zone_id and gateway_id

```hcl
data "tencentcloud_teo_multi_path_gateway_origin_acl" "example" {
  zone_id    = "zone-2noqxz9b6klw"
  gateway_id = "gw-12345678"
}
```
