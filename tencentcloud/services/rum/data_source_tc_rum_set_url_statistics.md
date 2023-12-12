Use this data source to query detailed information of rum set_url_statistics

Example Usage

```hcl
data "tencentcloud_rum_set_url_statistics" "set_url_statistics" {
  start_time = 1625444040
  type       = "allcount"
  end_time   = 1625454840
  project_id = 1
}
```