Provides a resource to create a CLS (Cloud Log Service) remote write task.

Example Usage

```hcl
resource "tencentcloud_cls_remote_write_task" "example" {
  topic_id             = "3f5cc19d-f601-48b6-9671-9019b7b09bd0"
  name                 = "tf-example"
  target               = "TencentCloud_Prometheus"
  remote_write_url     = "http://172.0.0.10:9090/api/v1/prom/write"
  auth_type            = 1
  net_type             = 1
  vpc_id               = "vpc-mkegckdp"
  virtual_gateway_type = 1025
  instance_id          = "prom-qha7cws8"
  has_services_log     = 2
  enable               = 1

  auth_info {
    username = "root"
    password = "Fqzg*******************55Ko"
  }
}
```

Import

CLS remote write task can be imported using the composite id, e.g. the format is `topic_id#task_id`

```
terraform import tencentcloud_cls_remote_write_task.example 3f5cc19d-f601-48b6-9671-9019b7b09bd0#e0e10aab-5977-49b3-a4dd-106db8821b6a
```
