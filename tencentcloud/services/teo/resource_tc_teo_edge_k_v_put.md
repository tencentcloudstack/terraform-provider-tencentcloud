Provides a resource to create a TEO Edge KV key-value pair attachment

Example Usage

```hcl
resource "tencentcloud_teo_edge_k_v_put" "example" {
  zone_id   = "zone-2o3h21ed2t68"
  namespace = "example-namespace"
  key       = "example-key"
  value     = "example-value"
}
```

With expiration

```hcl
resource "tencentcloud_teo_edge_k_v_put" "example_with_ttl" {
  zone_id        = "zone-2o3h21ed2t68"
  namespace      = "example-namespace"
  key            = "example-key-ttl"
  value          = "example-value"
  expiration_ttl = 3600
}
```

Import

TEO Edge KV put can be imported using the zoneId#namespace#key, e.g.

```
terraform import tencentcloud_teo_edge_k_v_put.example zone-2o3h21ed2t68#example-namespace#example-key
```
