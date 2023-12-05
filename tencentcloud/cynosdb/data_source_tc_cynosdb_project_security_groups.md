Use this data source to query detailed information of cynosdb project_security_groups

Example Usage

```hcl
data "tencentcloud_cynosdb_project_security_groups" "project_security_groups" {
  project_id = 1250480
  search_key = "自定义模版"
}
```