Provides a resource to create a waf web_shell

Example Usage

```hcl
resource "tencentcloud_waf_web_shell" "example" {
  domain = "demo.waf.com"
  status = 0
}
```

Import

waf web_shell can be imported using the id, e.g.

```
terraform import tencentcloud_waf_web_shell.example demo.waf.com
```