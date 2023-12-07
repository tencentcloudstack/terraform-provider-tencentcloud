Use this data source to query CDH instances.

Example Usage

```hcl
data "tencentcloud_cdh_instances" "list" {
  availability_zone = "ap-guangzhou-3"
  host_id = "host-d6s7i5q4"
  host_name = "test"
  host_state = "RUNNING"
  project_id = 1154137
}
```