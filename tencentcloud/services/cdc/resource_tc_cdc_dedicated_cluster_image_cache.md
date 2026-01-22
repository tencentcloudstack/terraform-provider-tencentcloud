Provides a resource to create a CDC dedicated cluster image cache

Example Usage

```hcl
resource "tencentcloud_cdc_dedicated_cluster_image_cache" "cdc_dedicated_cluster_image_cache" {
  dedicated_cluster_id = "cluster-262n63e8"
  image_id = "img-eb30mz89"
}
```

Import

CDC dedicated cluster image cache can be imported using the id, e.g.

```
terraform import tencentcloud_cdc_dedicated_cluster_image_cache.cdc_dedicated_cluster_image_cache ${dedicated_cluster_id}#${image_id}
```