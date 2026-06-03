Provides a resource to create a TEO edge KV namespace

Example Usage

```hcl
resource "tencentcloud_teo_edge_k_v_namespace" "example" {
  zone_id   = "zone-2o3h21ed2t68"
  namespace = "my-namespace"
  remark    = "This is an example namespace"
}
```

Import

TEO edge KV namespace can be imported using the composite id (zone_id#namespace), e.g.

````
terraform import tencentcloud_teo_edge_k_v_namespace.example zone-2o3h21ed2t68#my-namespace
````
