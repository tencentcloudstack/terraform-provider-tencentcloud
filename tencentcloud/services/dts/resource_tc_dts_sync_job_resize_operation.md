Provides a resource to create a dts sync_job_resize_operation

Example Usage

```hcl
resource "tencentcloud_dts_sync_job_resize_operation" "sync_job_resize_operation" {
  job_id = "sync-werwfs23"
  new_instance_class = "large"
}
```