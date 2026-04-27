Provides a resource to create a TEO security js injection rule

Example Usage

```hcl
resource "tencentcloud_teo_security_js_injection_rule" "example" {
  zone_id = "zone-2qtuhspy7cr6"
  js_injection_rules {
    name      = "test-rule-1"
    condition = "${http.request.host} in ['example.com']"
    inject_js = "inject-sdk-only"
  }
  js_injection_rules {
    name      = "test-rule-2"
    priority  = 10
    condition = "${http.request.uri.path} in ['/api/*']"
    inject_js = "no-injection"
  }
}
```

Import

TEO security js injection rule can be imported using the zone_id, e.g.

```
terraform import tencentcloud_teo_security_js_injection_rule.example zone-2qtuhspy7cr6
```
