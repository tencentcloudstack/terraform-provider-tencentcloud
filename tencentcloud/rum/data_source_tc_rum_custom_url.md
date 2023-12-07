Use this data source to query detailed information of rum custom_url

Example Usage

```hcl
data "tencentcloud_rum_custom_url" "custom_url" {
  start_time = 1625444040
  type = "top"
  end_time = 1625454840
  project_id = 1
}
```