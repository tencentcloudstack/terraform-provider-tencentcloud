Provides a resource to create a TEO (EdgeOne) alias domain.

Example Usage

Alias domain with no certificate

```hcl
resource "tencentcloud_teo_alias_domain" "example" {
  zone_id     = "zone-297z8rf93cfw"
  alias_name  = "alias.example.com"
  target_name = "target.example.com"
}
```

Alias domain with SSL managed certificate

```hcl
resource "tencentcloud_teo_alias_domain" "example" {
  zone_id     = "zone-297z8rf93cfw"
  alias_name  = "alias.example.com"
  target_name = "target.example.com"
  cert_type   = "hosting"
  cert_id     = ["cert-abc123"]
}
```

Import

teo alias_domain can be imported using the zone_id#alias_name, e.g.

```
terraform import tencentcloud_teo_alias_domain.example zone-297z8rf93cfw#alias.example.com
```
