Use this data source to query detailed information of rum static_url

Example Usage

```hcl
data "tencentcloud_rum_static_url" "static_url" {
  start_time = 1625444040
  type       = "pagepv"
  end_time   = 1625454840
  project_id = 1
}
```