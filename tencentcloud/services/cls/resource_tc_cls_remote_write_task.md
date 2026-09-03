Provides a resource to create a CLS (Cloud Log Service) remote write task.

Example Usage

```hcl
resource "tencentcloud_cls_remote_write_task" "example" {
  topic_id             = "d22f6119-68d7-4fce-abb1-4db8518b5ec4"
  name                 = "tf-example"
  target               = "TencentCloud_Prometheus"
  remote_write_url     = "http://172.16.0.14:9090/api/v1/prom/write"
  auth_type            = 1
  net_type             = 1
  vpc_id               = "vpc-mkegckdp"
  virtual_gateway_type = 1025
  instance_id          = "prom-qha7cws8"
  has_services_log     = 2

  auth_info {
    token    = "1309118522"
    password = "FqzgGX3Ty9TQs10ZtVaD5d255Ko"
  }
}
```

Import

CLS remote write task can be imported using the composite id, e.g. the format is `topic_id#task_id`

```
terraform import tencentcloud_cls_remote_write_task.example d22f6119-68d7-4fce-abb1-4db8518b5ec4#task_id
```
