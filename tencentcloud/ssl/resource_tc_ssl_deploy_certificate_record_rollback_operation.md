Provides a resource to create a ssl deploy_certificate_record_rollback

Example Usage

```hcl
resource "tencentcloud_ssl_deploy_certificate_record_rollback_operation" "deploy_certificate_record_rollback" {
  deploy_record_id = 35471
}
```

Import

ssl deploy_certificate_record_rollback can be imported using the id, e.g.

```
terraform import tencentcloud_ssl_deploy_certificate_record_rollback_operation.deploy_certificate_record_rollback deploy_certificate_record_rollback_id
```