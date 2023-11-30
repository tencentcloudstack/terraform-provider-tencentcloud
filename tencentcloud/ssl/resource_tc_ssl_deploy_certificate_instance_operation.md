Provides a resource to create a ssl deploy_certificate_instance

Example Usage

```hcl
resource "tencentcloud_ssl_deploy_certificate_instance_operation" "deploy_certificate_instance" {
  certificate_id = "8x1eUSSl"
  instance_id_list = ["cdndomain1.example.com|on","cdndomain1.example.com|off"]
  resource_type = "cdn"
}
```

Import

ssl deploy_certificate_instance can be imported using the id, e.g.

```
terraform import tencentcloud_ssl_deploy_certificate_instance_operation.deploy_certificate_instance deploy_certificate_instance_id
```