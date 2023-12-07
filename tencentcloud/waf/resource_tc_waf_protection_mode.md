Provides a resource to create a waf protection_mode

Example Usage

```hcl
resource "tencentcloud_waf_protection_mode" "example" {
  domain  = "keep.qcloudwaf.com"
  mode    = 10
  edition = "sparta-waf"
  type    = 0
}
```