Provides a resource to create a css callback_rule

Example Usage

```hcl
resource "tencentcloud_css_callback_rule_attachment" "callback_rule" {
  domain_name = "177154.push.tlivecloud.com"
  template_id = 434039
  app_name    = "live"
}
```

Import

css callback_rule can be imported using the id, e.g.

```
terraform import tencentcloud_css_callback_rule_attachment.callback_rule templateId#domainName
```