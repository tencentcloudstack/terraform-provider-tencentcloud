Provides a resource to create a ssl update_certificate_record_retry

Example Usage

```hcl
resource "tencentcloud_ssl_update_certificate_record_retry_operation" "update_certificate_record_retry" {
  deploy_record_id = "1603"
}
```

Import

ssl update_certificate_record_retry can be imported using the id, e.g.

```
terraform import tencentcloud_ssl_update_certificate_record_retry_operation.update_certificate_record_retry update_certificate_record_retry_id
```