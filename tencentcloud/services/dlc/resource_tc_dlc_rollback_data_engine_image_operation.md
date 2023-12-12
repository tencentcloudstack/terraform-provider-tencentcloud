Provides a resource to create a dlc rollback_data_engine_image

Example Usage

```hcl
data "tencentcloud_dlc_check_data_engine_image_can_be_rollback" "check_data_engine_image_can_be_rollback" {
  data_engine_id = "DataEngine-cgkvbas6"
}
resource "tencentcloud_dlc_rollback_data_engine_image_operation" "rollback_data_engine_image" {
  data_engine_id = "DataEngine-cgkvbas6"
  from_record_id = data.tencentcloud_dlc_check_data_engine_image_can_be_rollback.check_data_engine_image_can_be_rollback.from_record_id
  to_record_id = data.tencentcloud_dlc_check_data_engine_image_can_be_rollback.check_data_engine_image_can_be_rollback.to_record_id
}
```

Import

dlc rollback_data_engine_image can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_rollback_data_engine_image_operation.rollback_data_engine_image rollback_data_engine_image_id
```