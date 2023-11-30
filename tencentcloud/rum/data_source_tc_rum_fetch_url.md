Use this data source to query detailed information of rum fetch_url

Example Usage

```hcl
data "tencentcloud_rum_fetch_url" "fetch_url" {
  start_time = 1625444040
  type = "allcount"
  end_time = 1625454840
  project_id = 1
}
```