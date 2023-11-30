Provides a resource to create a mariadb switch_h_a

Example Usage

```hcl
resource "tencentcloud_mariadb_switch_ha" "switch_ha" {
  instance_id = "tdsql-9vqvls95"
  zone        = "ap-guangzhou-2"
}
```