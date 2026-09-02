Provides a resource to create a CLS (Cloud Log Service) remote write task.

Example Usage

```hcl
resource "tencentcloud_cls_logset" "example" {
  logset_name = "tf-example-remote-write"
  tags = {
    createBy = "Terraform"
  }
}

resource "tencentcloud_cls_topic" "example" {
  topic_name           = "tf-example-remote-write"
  logset_id            = tencentcloud_cls_logset.example.id
  auto_split           = false
  max_split_partitions = 20
  partition_count      = 1
  period               = 10
  storage_type         = "hot"
  tags = {
    createBy = "Terraform"
  }
}

resource "tencentcloud_cls_remote_write_task" "example" {
  topic_id         = tencentcloud_cls_topic.example.id
  name             = "tf-example-remote-write"
  target           = "prometheus"
  remote_write_url = "http://prometheus-server.monitoring.svc.cluster.local:9090/api/v1/write"
  auth_type        = 0
  net_type         = 1
  vpc_id           = "vpc-xxxxxxxx"
  enable           = 1
}
```

Example with basic auth:

```hcl
resource "tencentcloud_cls_remote_write_task" "basic_auth_example" {
  topic_id         = tencentcloud_cls_topic.example.id
  name             = "tf-example-remote-write-basic-auth"
  target           = "prometheus"
  remote_write_url = "https://prometheus.example.com/api/v1/write"
  auth_type        = 1
  net_type         = 2

  auth_info {
    username = "admin"
    password = "my-password"
  }
}
```

Example with token auth:

```hcl
resource "tencentcloud_cls_remote_write_task" "token_auth_example" {
  topic_id         = tencentcloud_cls_topic.example.id
  name             = "tf-example-remote-write-token-auth"
  target           = "prometheus"
  remote_write_url = "https://prometheus.example.com/api/v1/write"
  auth_type        = 2
  net_type         = 2

  auth_info {
    token = "my-token-string"
  }
}
```

Import

CLS remote write task can be imported using the composite id, e.g. the format is `task_id#topic_id`

```
$ terraform import tencentcloud_cls_remote_write_task.example task-id-xxxxxx#topic-id-xxxxxx
```
