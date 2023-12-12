Provides a resource to create a css record_rule

Example Usage

```hcl
resource "tencentcloud_css_record_rule_attachment" "record_rule" {
  domain_name = "177154.push.tlivecloud.com"
  template_id = 1262818
  app_name    = "qqq"
  stream_name = "ppp"
}
```

Import

css record_rule can be imported using the id, e.g.

```
terraform import tencentcloud_css_record_rule_attachment.record_rule templateId#domainName
```