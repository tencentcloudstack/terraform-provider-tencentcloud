Provides a resource to create a TEO prefetch origin limit config.

Example Usage

Set prefetch origin bandwidth limit for overseas area

```hcl
resource "tencentcloud_teo_prefetch_origin_limit" "example" {
  zone_id     = "zone-3edjdliiw3he"
  domain_name = "example.com"
  area        = "Overseas"
  bandwidth   = 200
}
```

Set prefetch origin bandwidth limit for Mainland China area

```hcl
resource "tencentcloud_teo_prefetch_origin_limit" "example" {
  zone_id     = "zone-3edjdliiw3he"
  domain_name = "example.com"
  area        = "MainlandChina"
  bandwidth   = 500
}
```

Import

TEO prefetch origin limit config can be imported using the composite ID format `zone_id#domain_name#area`, e.g.

```
terraform import tencentcloud_teo_prefetch_origin_limit.example zone-3edjdliiw3he#example.com#Overseas
```
