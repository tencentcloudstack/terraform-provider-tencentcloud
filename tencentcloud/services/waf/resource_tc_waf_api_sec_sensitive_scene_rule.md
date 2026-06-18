Provides a resource to create a WAF api sec sensitive scene rule

Example Usage

```hcl
resource "tencentcloud_waf_api_sec_sensitive_scene_rule" "example" {
  domain    = "www.example.com"
  rule_name = "tf-example"
  status    = 1
  source    = "custom"

  rule_list {
    key     = "api"
    operate = "equal"
    value   = ["/login", "/user"]
  }
}
```

Import

WAF api sec sensitive scene rule can be imported using the domain#ruleName, e.g.

```
terraform import tencentcloud_waf_api_sec_sensitive_scene_rule.example www.example.com#tf-example
```
