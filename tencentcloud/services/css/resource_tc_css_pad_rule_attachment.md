Provides a resource to create a css pad_rule_attachment

Example Usage

```hcl
resource "tencentcloud_css_pad_rule_attachment" "pad_rule_attachment" {
  domain_name = "177154.push.tlivecloud.com"
  template_id = 17067
  app_name    = "qqq"
  stream_name = "ppp"
}
```

Import

css pad_rule_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_css_pad_rule_attachment.pad_rule_attachment templateId#domainName
```