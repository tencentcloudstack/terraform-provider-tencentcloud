Provides a resource to create a teo origin_group

~> **NOTE:** Please note that `tencentcloud_teo_origin_group` had to undergo incompatible changes in version v1.81.96.

Example Usage

Self origin group

```hcl
resource "tencentcloud_teo_origin_group" "basic" {
  name    = "keep-group-1"
  type    = "GENERAL"
  zone_id = "zone-197z8rf93cfw"

  records {
    record  = "tf-teo.xyz"
    type    = "IP_DOMAIN"
    weight  = 100
    private = true

    private_parameters {
      name = "SecretAccessKey"
      value = "test"
    }
  }
}
```
Import

teo origin_group can be imported using the zone_id#originGroup_id, e.g.
````
terraform import tencentcloud_teo_origin_group.origin_group zone-297z8rf93cfw#origin-4f8a30b2-3720-11ed-b66b-525400dceb86
````