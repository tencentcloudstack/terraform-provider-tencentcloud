Use this data source to query detailed information of bi project

Example Usage

```hcl
data "tencentcloud_bi_project" "project" {
  page_no = 1
  keyword = "abc"
  all_page = true
  module_collection = "sys_common_user"
}
```