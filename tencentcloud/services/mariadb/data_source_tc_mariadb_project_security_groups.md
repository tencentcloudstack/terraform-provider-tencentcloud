Use this data source to query detailed information of mariadb project_security_groups

Example Usage

```hcl
data "tencentcloud_mariadb_project_security_groups" "project_security_groups" {
  product    = "mariadb"
  project_id = 0
}
```