Provides a resource to create a teo zone

Example Usage

```hcl
resource "tencentcloud_teo_zone" "zone" {
  zone_name       = "tf-teo.com"
  type            = "partial"
  area            = "overseas"
  alias_zone_name = "teo-test"
  paused          = false
  plan_id         = "edgeone-2kfv1h391n6w"
  tags = {
    "createdBy" = "terraform"
  }
}

```
Import

teo zone can be imported using the id, e.g.
```
terraform import tencentcloud_teo_zone.zone zone_id
```