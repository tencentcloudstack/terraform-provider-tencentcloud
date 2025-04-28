Provides a resource to create a cls cloud product log task

~> **NOTE:** This resource has been deprecated in Terraform TencentCloud provider version 1.81.188. Please use `tencentcloud_cls_cloud_product_log_task_v2` instead.

~> **NOTE:** Using this resource will create new `logset` and `topic`

Example Usage

```hcl
resource "tencentcloud_cls_cloud_product_log_task" "example" {
  instance_id          = "postgres-1p7xvpc1"
  assumer_name         = "PostgreSQL"
  log_type             = "PostgreSQL-SLOW"
  cloud_product_region = "gz"
  cls_region           = "ap-guangzhou"
  logset_name          = "tf-example"
  topic_name           = "tf-example"
}
```

Import

cls cloud product log task can be imported using the id, e.g.

```
terraform import tencentcloud_cls_cloud_product_log_task.example postgres-1p7xvpc1#PostgreSQL#PostgreSQL-SLOW#gz
```
