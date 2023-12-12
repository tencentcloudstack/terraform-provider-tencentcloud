Provides a resource to create a css backup_stream

~> **NOTE:** This resource is only valid when the push stream. When the push stream ends, it will be deleted.

Example Usage

```hcl
resource "tencentcloud_css_backup_stream" "backup_stream" {
  push_domain_name  = "177154.push.tlivecloud.com"
  app_name          = "live"
  stream_name       = "1308919341_test"
  upstream_sequence = "2209501773993286139"
}
```

Import

css backup_stream can be imported using the id, e.g.

```
terraform import tencentcloud_css_backup_stream.backup_stream pushDomainName#appName#streamName
```