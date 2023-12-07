Provides a resource to create a dlc renew_data_engine

Example Usage

```hcl
resource "tencentcloud_dlc_renew_data_engine_operation" "renew_data_engine" {
  data_engine_name = "testEngine"
  time_span = 3600
  pay_mode = 1
  time_unit = "m"
  renew_flag = 1
}
```

Import

dlc renew_data_engine can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_renew_data_engine_operation.renew_data_engine renew_data_engine_id
```