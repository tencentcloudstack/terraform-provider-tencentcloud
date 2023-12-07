Provides a resource to create a ssl deploy_certificate_record_retry

Example Usage

```hcl
resource "tencentcloud_ssl_deploy_certificate_record_retry_operation" "deploy_certificate_record_retry" {
  deploy_record_id = 35474
}
```

Import

ssl deploy_certificate_record_retry can be imported using the id, e.g.

```
terraform import tencentcloud_ssl_deploy_certificate_record_retry_operation.deploy_certificate_record_retry deploy_certificate_record_retry_id
```