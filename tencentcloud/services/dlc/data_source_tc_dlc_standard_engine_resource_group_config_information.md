Use this data source to query detailed information of DLC standard engine resource group config information

Example Usage

```hcl
data "tencentcloud_dlc_standard_engine_resource_group_config_information" "example" {
  sort_by = "create-time"
  sorting = "desc"
  filters {
    name = "engine-id"
    values = [
      "DataEngine-5plqp7q7"
    ]
  }
}
```
