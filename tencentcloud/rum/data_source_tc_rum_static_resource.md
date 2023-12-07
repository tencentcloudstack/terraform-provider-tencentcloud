Use this data source to query detailed information of rum static_resource

Example Usage

```hcl
data "tencentcloud_rum_static_resource" "static_resource" {
  start_time = 1625444040
  type       = "top"
  end_time   = 1625454840
  project_id = 1
}
```