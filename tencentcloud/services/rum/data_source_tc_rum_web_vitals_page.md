Use this data source to query detailed information of rum web_vitals_page

Example Usage

```hcl
data "tencentcloud_rum_web_vitals_page" "web_vitals_page" {
  start_time = 1625444040
  end_time   = 1625454840
  project_id = 1
  type       = "from"
}
```