Use this data source to query detailed information of mariadb upgrade_price

Example Usage

```hcl
data "tencentcloud_mariadb_upgrade_price" "upgrade_price" {
  instance_id = "tdsql-9vqvls95"
  memory      = 4
  storage     = 40
  node_count  = 2
}
```