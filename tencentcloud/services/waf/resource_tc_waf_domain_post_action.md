Provides a resource to create a WAF domain post action

Example Usage

```hcl
resource "tencentcloud_waf_domain_post_action" "example" {
  domain             = "example.com"
  post_cls_action    = 1
  post_ckafka_action = 0
}
```

Import

WAF domain post action can be imported using the id, e.g.

```
terraform import tencentcloud_waf_domain_post_action.example example.com
```
