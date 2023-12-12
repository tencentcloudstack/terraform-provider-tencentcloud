Use this data source to query detailed information of rum fetch_url_info

Example Usage

```hcl
data "tencentcloud_rum_fetch_url_info" "fetch_url_info" {
  start_time = 1625444040
  type = "top"
  end_time = 1625454840
  project_id = 1
}
```