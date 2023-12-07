Provides a resource to create a oceanus job

Example Usage

```hcl
resource "tencentcloud_oceanus_job" "example" {
  name          = "tf_example_job"
  job_type      = 1
  cluster_type  = 2
  cluster_id    = "cluster-1kcd524h"
  cu_mem        = 4
  remark        = "remark."
  folder_id     = "folder-7ctl246z"
  flink_version = "Flink-1.16"
  work_space_id = "space-2idq8wbr"
}
```