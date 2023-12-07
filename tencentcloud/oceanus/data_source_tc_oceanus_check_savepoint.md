Use this data source to query detailed information of oceanus check_savepoint

Example Usage

```hcl
data "tencentcloud_oceanus_check_savepoint" "example" {
  job_id         = "cql-314rw6w0"
  serial_id      = "svp-52xkpymp"
  record_type    = 1
  savepoint_path = "cosn://52xkpymp-12345/12345/10000/cql-12345/2/flink-savepoints/savepoint-000000-12334"
  work_space_id  = "space-2idq8wbr"
}
```