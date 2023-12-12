Provides a resource to create a oceanus job_copy

Example Usage

```hcl
resource "tencentcloud_oceanus_job_copy" "example" {
  source_id         = "cql-0nob2hx8"
  target_cluster_id = "cluster-1kcd524h"
  source_name       = "keep_jar"
  target_name       = "tf_copy_example"
  target_folder_id  = "folder-7ctl246z"
  job_type          = 2
  work_space_id     = "space-2idq8wbr"
}
```