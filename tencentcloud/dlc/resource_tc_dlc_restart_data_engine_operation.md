Provides a resource to create a dlc restart_data_engine

Example Usage

```hcl
resource "tencentcloud_dlc_restart_data_engine_operation" "restart_data_engine" {
  data_engine_id = "DataEngine-g5ds87d8"
  forced_operation = false
}
```