Use this data source to query detailed information of mariadb databases

Example Usage

```hcl
data "tencentcloud_mariadb_databases" "databases" {
  instance_id = "tdsql-e9tklsgz"
}
```