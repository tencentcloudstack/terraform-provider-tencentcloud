Provides a resource to apply TEO (EdgeOne) free certificate for a domain.

Example Usage

Apply free certificate with DNS verification

```hcl
resource "tencentcloud_teo_apply_free_certificate" "example" {
  zone_id             = "zone-2o3h21ed8bsf"
  domain              = "www.example.com"
  verification_method = "dns_challenge"
}
```

Apply free certificate with HTTP file verification

```hcl
resource "tencentcloud_teo_apply_free_certificate" "example_http" {
  zone_id             = "zone-2o3h21ed8bsf"
  domain              = "www.example.com"
  verification_method = "http_challenge"
}
```
