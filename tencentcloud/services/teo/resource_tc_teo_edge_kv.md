Provides a resource to create a TEO Edge KV key-value pair

Example Usage

```hcl
resource "tencentcloud_teo_edge_kv" "example" {
  zone_id   = "zone-2o3h21ed2t68"
  namespace = "example-namespace"
  key       = "example-key"
  value     = "example-value"
}
```

Import

TEO Edge KV can be imported using the zoneId#namespace#key, e.g.

```
terraform import tencentcloud_teo_edge_kv.example zone-2o3h21ed2t68#example-namespace#example-key
```
