Use this data source to query detailed information of oceanus work_spaces

Example Usage

```hcl
data "tencentcloud_oceanus_work_spaces" "example" {
  order_type = 1
  filters {
    name   = "WorkSpaceName"
    values = ["tf_example"]
  }
}
```