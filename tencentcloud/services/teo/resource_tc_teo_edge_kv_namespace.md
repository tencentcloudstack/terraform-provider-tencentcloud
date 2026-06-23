Provides a resource to create a TEO Edge KV namespace

Example Usage

```hcl
resource "tencentcloud_teo_edge_kv_namespace" "example" {
  zone_id   = "zone-2o3h21ed2t68"
  namespace = "example-namespace"
  remark    = "This is an example namespace"
}
```

Import

TEO Edge KV namespace can be imported using the zone_id#namespace, e.g.

```
terraform import tencentcloud_teo_edge_kv_namespace.example zone-2o3h21ed2t68#example-namespace
```
