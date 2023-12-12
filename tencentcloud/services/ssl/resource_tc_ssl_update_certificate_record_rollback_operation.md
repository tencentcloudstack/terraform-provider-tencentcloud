Provides a resource to create a ssl update_certificate_record_rollback

Example Usage

```hcl
resource "tencentcloud_ssl_update_certificate_record_rollback_operation" "update_certificate_record_rollback" {
  deploy_record_id = "1603"
}
```

Import

ssl update_certificate_record_rollback can be imported using the id, e.g.

```
terraform import tencentcloud_ssl_update_certificate_record_rollback_operation.update_certificate_record_rollback update_certificate_record_rollback_id
```