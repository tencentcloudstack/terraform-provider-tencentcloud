Use this data source to query detailed information of rum static_project

Example Usage

```hcl
data "tencentcloud_rum_static_project" "static_project" {
  start_time = 1625444040
  type       = "allcount"
  end_time   = 1625454840
  project_id = 1
}
```