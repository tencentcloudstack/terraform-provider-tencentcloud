Provides a resource to create a DLC rollback data engine image

Example Usage

```hcl
data "tencentcloud_dlc_check_data_engine_image_can_be_rollback" "example" {
  data_engine_id = "DataEngine-cgkvbas6"
}

resource "tencentcloud_dlc_rollback_data_engine_image_operation" "example" {
  data_engine_id = "DataEngine-cgkvbas6"
  from_record_id = data.tencentcloud_dlc_check_data_engine_image_can_be_rollback.example.from_record_id
  to_record_id   = data.tencentcloud_dlc_check_data_engine_image_can_be_rollback.example.to_record_id
}
```
