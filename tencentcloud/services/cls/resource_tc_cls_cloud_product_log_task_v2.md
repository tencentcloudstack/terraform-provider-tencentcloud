Provides a resource to create a cls cloud product log task

~> **NOTE:** In the destruction of resources, if cascading deletion of logset and topic is required, please set `force_delete` to `true`.

Example Usage

Create log delivery using the default newly created logset and topic

```hcl
resource "tencentcloud_cls_cloud_product_log_task_v2" "example" {
  instance_id          = "postgres-0an6hpv3"
  assumer_name         = "PostgreSQL"
  log_type             = "PostgreSQL-SLOW"
  cloud_product_region = "gz"
  cls_region           = "ap-guangzhou"
  logset_name          = "tf-example"
  topic_name           = "tf-example"
  force_delete         = true
}
```

Create log delivery using existing logset and topic

```hcl
resource "tencentcloud_cls_cloud_product_log_task_v2" "example" {
  instance_id          = "postgres-0an6hpv3"
  assumer_name         = "PostgreSQL"
  log_type             = "PostgreSQL-SLOW"
  cloud_product_region = "gz"
  cls_region           = "ap-guangzhou"
  logset_id            = "ca5b4f56-1174-4eee-bc4c-69e48e0e8c45"
  topic_id             = "d8177ca9-466b-42f4-a110-5933daf0a83a"
  force_delete         = false
}
```

Import

cls cloud product log task can be imported using the id, e.g.

```
terraform import tencentcloud_cls_cloud_product_log_task_v2.example postgres-1p7xvpc1#PostgreSQL#PostgreSQL-SLOW#gz
```
