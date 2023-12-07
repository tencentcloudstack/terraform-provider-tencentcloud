Provides a resource to create a waf cc_auto_status

Example Usage

```hcl
resource "tencentcloud_waf_cc_auto_status" "example" {
  domain  = "www.demo.com"
  edition = "sparta-waf"
}
```

Import

waf cc_auto_status can be imported using the id, e.g.

```
terraform import tencentcloud_waf_cc_auto_status.example www.demo.com#sparta-waf
```