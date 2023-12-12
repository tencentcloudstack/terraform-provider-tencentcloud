Provides a resource to create a mariadb renew_instance

Example Usage

```hcl
resource "tencentcloud_mariadb_renew_instance" "renew_instance" {
  instance_id = "tdsql-9vqvls95"
  period      = 1
}
```