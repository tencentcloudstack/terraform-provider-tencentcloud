Provides a resource to create a ssl update_certificate_instance

Example Usage

```hcl
resource "tencentcloud_ssl_update_certificate_instance_operation" "update_certificate_instance" {
  certificate_id = "8x1eUSSl"
  old_certificate_id = "8xNdi2ig"
  resource_types = ["cdn"]
}
```
Upload certificate

```hcl
resource "tencentcloud_ssl_update_certificate_instance_operation" "update_certificate_instance" {
  old_certificate_id = "xxx"
  certificate_public_key = file("xxx.crt")
  certificate_private_key= file("xxx.key")
  repeatable= true
  resource_types = ["cdn"]
}
```
