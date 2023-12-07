Use this data source to query detailed information of wedata data_source_without_info

Example Usage

```hcl
data "tencentcloud_wedata_data_source_without_info" "example" {
  filters {
    name   = "ownerProjectId"
    values = ["1612982498218618880"]
  }

  order_fields {
    name      = "create_time"
    direction = "DESC"
  }
}
```