Use this data source to query detailed information of wedata data_source_list

Example Usage

Query All

```hcl
data "tencentcloud_wedata_data_source_list" "example" {}
```

Query By filter

```hcl
data "tencentcloud_wedata_data_source_list" "example" {
  order_fields {
    name      = "create_time"
    direction = "DESC"
  }

  filters {
    name   = "Name"
    values = ["tf_example"]
  }
}
```