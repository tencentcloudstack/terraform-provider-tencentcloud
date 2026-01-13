Provides a resource to create a IGTM package task

Example Usage

```hcl
resource "tencentcloud_igtm_package_task" "example" {
  task_detection_quantity = 100
  auto_renew              = 2
  time_span               = 1
  auto_voucher            = 1
}
```

Import

IGTM package task can be imported using the id, e.g.

```
terraform import tencentcloud_igtm_package_task.example task-dahygvmzawgn
```
