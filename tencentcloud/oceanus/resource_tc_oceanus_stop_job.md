Provides a resource to create a oceanus stop_job

Example Usage

```hcl
resource "tencentcloud_oceanus_stop_job" "example" {
  stop_job_descriptions {
    job_id    = "cql-4xwincyn"
    stop_type = 1
  }
  work_space_id = "space-2idq8wbr"
}
```