Provides a resource to create a mariadb encrypt_attributes

Example Usage

```hcl
resource "tencentcloud_mariadb_encrypt_attributes" "encrypt_attributes" {
  instance_id = "tdsql-ow728lmc"
  encrypt_enabled = 1
}
```
