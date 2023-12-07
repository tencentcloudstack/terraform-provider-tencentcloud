Provides a resource to create a waf module_status

Example Usage

```hcl
resource "tencentcloud_waf_module_status" "example" {
  domain         = "demo.waf.com"
  web_security   = 1
  access_control = 0
  cc_protection  = 1
  api_protection = 1
  anti_tamper    = 1
  anti_leakage   = 0
}
```

Import

waf module_status can be imported using the id, e.g.

```
terraform import tencentcloud_waf_module_status.example demo.waf.com
```