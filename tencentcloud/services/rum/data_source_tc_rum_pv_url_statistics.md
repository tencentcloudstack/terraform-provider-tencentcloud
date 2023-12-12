Use this data source to query detailed information of rum pv_url_statistics

Example Usage

```hcl
data "tencentcloud_rum_pv_url_statistics" "pv_url_statistics" {
  start_time = 1625444040
  type       = "allcount"
  end_time   = 1625454840
  project_id = 1
}
```