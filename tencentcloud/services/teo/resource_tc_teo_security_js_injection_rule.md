Provides a resource to create a TEO security JavaScript injection rule.

Example Usage

```hcl
resource "tencentcloud_teo_security_js_injection_rule" "example" {
  zone_id = "zone-3fkff38fyw8s"

  js_injection_rules {
    name      = "tf-example"
    priority  = 50
    condition = "$${http.request.host} in ['www.demo.com']"
    inject_j_s = "inject-sdk-only"
  }
}
```

Import

TEO security JavaScript injection rule can be imported using the zoneId#ruleId, e.g.

```
terraform import tencentcloud_teo_security_js_injection_rule.example zone-3fkff38fyw8s#inject-0000040467
```
