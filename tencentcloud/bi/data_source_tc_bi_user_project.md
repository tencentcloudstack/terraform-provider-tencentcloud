Use this data source to query detailed information of bi user_project

Example Usage

```hcl
data "tencentcloud_bi_user_project" "user_project" {
  project_id = 123
  all_page = true
}
```