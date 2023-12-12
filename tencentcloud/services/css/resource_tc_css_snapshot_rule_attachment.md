Provides a resource to create a css snapshot_rule

Example Usage

```hcl
resource "tencentcloud_css_snapshot_rule_attachment" "snapshot_rule" {
  domain_name = "177154.push.tlivecloud.com"
  template_id = 12838073
  app_name    = "qqq"
  stream_name = "ppp"
}
```

Import

css snapshot_rule can be imported using the id, e.g.

```
terraform import tencentcloud_css_snapshot_rule_attachment.snapshot_rule templateId#domainName
```