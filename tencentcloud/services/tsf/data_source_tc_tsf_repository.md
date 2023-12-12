Use this data source to query detailed information of tsf repository

Example Usage

```hcl
data "tencentcloud_tsf_repository" "repository" {
  search_word = "test"
  repository_type = "default"
}
```