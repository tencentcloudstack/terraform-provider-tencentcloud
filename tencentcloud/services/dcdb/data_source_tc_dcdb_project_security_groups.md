Use this data source to query detailed information of dcdb project_security_groups

Example Usage

```hcl
data "tencentcloud_dcdb_project_security_groups" "project_security_groups" {
  product    = "dcdb"
  project_id = 0
}
```