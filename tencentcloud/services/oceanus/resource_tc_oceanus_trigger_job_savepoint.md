Provides a resource to create a oceanus trigger_job_savepoint

Example Usage

```hcl
resource "tencentcloud_oceanus_trigger_job_savepoint" "example" {
  job_id        = "cql-4xwincyn"
  description   = "description."
  work_space_id = "space-2idq8wbr"
}
```