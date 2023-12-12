Provides a resource to create a teo origin_group

Example Usage

Self origin group

```hcl
resource "tencentcloud_teo_origin_group" "origin_group" {
  zone_id            = "zone-297z8rf93cfw"
  configuration_type = "weight"
  origin_group_name  = "test-group"
  origin_type        = "self"
  origin_records {
    area    = []
    port    = 8080
    private = false
    record  = "150.109.8.1"
    weight  = 100
  }
}

```

Cos origin group

```hcl
resource "tencentcloud_teo_origin_group" "origin_group" {
  configuration_type = "weight"
  origin_group_name  = "test"
  origin_type        = "cos"
  zone_id            = "zone-2o3h21ed8bpu"

  origin_records {
    area    = []
    port    = 0
    private = true
    record  = "test-ruichaolin-1310708577.cos.ap-nanjing.myqcloud.com"
    weight  = 100
  }
}
```
Import

teo origin_group can be imported using the zone_id#originGroup_id, e.g.
````
terraform import tencentcloud_teo_origin_group.origin_group zone-297z8rf93cfw#origin-4f8a30b2-3720-11ed-b66b-525400dceb86
````