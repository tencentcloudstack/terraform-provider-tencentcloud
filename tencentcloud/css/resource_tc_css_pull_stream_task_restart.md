Provides a resource to create a css restart_push_task

Example Usage

```hcl
resource "tencentcloud_css_pull_stream_task_restart" "restart_push_task" {
  task_id  = "3573"
  operator = "tf-test"
}
```