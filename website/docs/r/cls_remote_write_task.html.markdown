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

CLS remote write task can be imported using the composite id, e.g. the format is `topic_id#task_id`

```
terraform import tencentcloud_cls_remote_write_task.example d22f6119-68d7-4fce-abb1-4db8518b5ec4#task_id
```

