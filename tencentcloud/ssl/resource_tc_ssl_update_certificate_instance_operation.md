Provides a resource to create a ssl update_certificate_instance

Example Usage

```hcl
resource "tencentcloud_ssl_update_certificate_instance_operation" "update_certificate_instance" {
  certificate_id = "8x1eUSSl"
  old_certificate_id = "8xNdi2ig"
  resource_types = ["cdn"]
}
```

Import

ssl update_certificate_instance can be imported using the id, e.g.

```
terraform import tencentcloud_ssl_update_certificate_instance_operation.update_certificate_instance update_certificate_instance_id
```