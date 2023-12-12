Use this data source to query detailed information of mariadb securityGroups

Example Usage

```hcl
data "tencentcloud_mariadb_security_groups" "securityGroups" {
  instance_id = "tdsql-4pzs5b67"
  product = "mariadb"
}
```