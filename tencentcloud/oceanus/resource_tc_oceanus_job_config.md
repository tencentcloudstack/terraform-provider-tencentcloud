Provides a resource to create a oceanus job_config

Example Usage

If `log_collect_type` is 2

```hcl
resource "tencentcloud_oceanus_job_config" "example" {
  job_id           = "cql-4xwincyn"
  entrypoint_class = "tf_example"
  program_args     = "--conf Key=Value"
  remark           = "remark."
  resource_refs {
    resource_id = "resource-q22ntswy"
    version     = 1
    type        = 1
  }
  default_parallelism = 1
  properties {
    key   = "pipeline.max-parallelism"
    value = "2048"
  }
  log_collect       = true
  job_manager_spec  = "1"
  task_manager_spec = "1"
  cls_logset_id     = "cd9adbb5-6b7d-48d2-9870-77658959c7a4"
  cls_topic_id      = "cec4c2f1-0bf3-470e-b1a5-b1c451e88838"
  log_collect_type  = 2
  work_space_id     = "space-2idq8wbr"
  log_level         = "INFO"
  auto_recover      = 1
  expert_mode_on    = false
}
```

If `log_collect_type` is 3

```hcl
resource "tencentcloud_oceanus_job_config" "example" {
  job_id           = "cql-4xwincyn"
  entrypoint_class = "tf_example"
  program_args     = "--conf Key=Value"
  remark           = "remark."
  resource_refs {
    resource_id = "resource-q22ntswy"
    version     = 1
    type        = 1
  }
  default_parallelism = 1
  properties {
    key   = "pipeline.max-parallelism"
    value = "2048"
  }
  log_collect       = true
  job_manager_spec  = "1"
  task_manager_spec = "1"
  cls_logset_id     = "cd9adbb5-6b7d-48d2-9870-77658959c7a4"
  cls_topic_id      = "cec4c2f1-0bf3-470e-b1a5-b1c451e88838"
  log_collect_type  = 3
  work_space_id     = "space-2idq8wbr"
  log_level         = "INFO"
  auto_recover      = 1
  expert_mode_on    = false
  cos_bucket        = "autotest-gz-bucket-1257058945"
}
```