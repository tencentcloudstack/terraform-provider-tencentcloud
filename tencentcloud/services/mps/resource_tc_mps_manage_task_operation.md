Provides a resource to create a mps manage_task_operation

Example Usage

```hcl
resource "tencentcloud_mps_manage_task_operation" "operation" {
  operation_type = "Abort"
  task_id = "2600010949-LiveScheduleTask-xxxx"
}
```