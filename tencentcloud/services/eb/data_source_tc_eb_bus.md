Use this data source to query detailed information of eb bus

Example Usage

```hcl
data "tencentcloud_eb_bus" "this" {
  order_by = "created_at"
  order    = "DESC"

  filters {
    name   = "Type"
    values = ["Cloud", "Platform"]
  }

  filters {
    name   = "EventBusName"
    values = ["default"]
  }
}
```