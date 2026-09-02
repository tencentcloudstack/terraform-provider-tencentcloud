---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_remote_write_task"
sidebar_current: "docs-tencentcloud-resource-cls_remote_write_task"
description: |-
  Provides a resource to create a CLS (Cloud Log Service) remote write task.
---

# tencentcloud_cls_remote_write_task

Provides a resource to create a CLS (Cloud Log Service) remote write task.

## Example Usage

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

### Example with basic auth:

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

### Example with token auth:

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

## Argument Reference

The following arguments are supported:

* `auth_type` - (Required, Int) Authentication type. 0: no auth, 1: basic_auth, 2: token.
* `name` - (Required, String) RemoteWrite task name.
* `net_type` - (Required, Int) Network type. 1: intranet, 2: internet.
* `remote_write_url` - (Required, String) Target address for RemoteWrite.
* `target` - (Required, String) Target service name.
* `topic_id` - (Required, String) Log topic ID.
* `auth_info` - (Optional, List) Authentication information block.
* `enable` - (Optional, Int) Task status. 0: disabled, 1: enabled.
* `virtual_gateway_type` - (Optional, Int) Backend service type. 0: CVM, 1025: CLB.
* `vpc_id` - (Optional, String) Private network ID.

The `auth_info` object supports the following:

* `password` - (Optional, String) Basic auth password.
* `token` - (Optional, String) Basic auth token.
* `username` - (Optional, String) Basic auth username.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Task creation time.
* `logset_id` - Logset ID.
* `status` - Task running status. 1: running, 2: paused, 3: failed.
* `task_id` - RemoteWrite task ID.
* `update_time` - Task update time.


## Import

CLS remote write task can be imported using the composite id, e.g. the format is `task_id#topic_id`

```
$ terraform import tencentcloud_cls_remote_write_task.example task-id-xxxxxx#topic-id-xxxxxx
```

