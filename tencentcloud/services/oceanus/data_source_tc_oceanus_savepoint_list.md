Use this data source to query detailed information of oceanus savepoint_list

Example Usage

```hcl
data "tencentcloud_oceanus_savepoint_list" "example" {
  job_id        = "cql-314rw6w0"
  work_space_id = "space-2idq8wbr"
}
```