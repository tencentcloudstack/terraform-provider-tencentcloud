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
* `has_services_log` - (Optional, Int) Enable delivery service logging. 1: Disabled, 2: Enabled. Default value: 2.
* `instance_id` - (Optional, String) Instance ID.
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

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateRemoteWriteTask](https://www.tencentcloud.com/document/product/614/81463) - This API is used to create a remote write delivery task.
* Create: [DescribeRemoteWriteTasks](https://www.tencentcloud.com/document/product/614/81461) - This API is used to obtain the list of RemoteWrite delivery tasks.
* Read: [DescribeRemoteWriteTasks](https://www.tencentcloud.com/document/product/614/81461) - This API is used to obtain the list of RemoteWrite delivery tasks.
* Update: [DescribeRemoteWriteTasks](https://www.tencentcloud.com/document/product/614/81461) - This API is used to obtain the list of RemoteWrite delivery tasks.
* Update: [ModifyRemoteWriteTask](https://www.tencentcloud.com/document/product/614/81460) - This API is used to modify a RemoteWrite task.
* Delete: [DeleteRemoteWriteTask](https://www.tencentcloud.com/document/product/614/81462) - This API is used to delete a RemoteWrite task.

## Import

CLS remote write task can be imported using the composite id, e.g. the format is `topic_id#task_id`

```
terraform import tencentcloud_cls_remote_write_task.example 3f5cc19d-f601-48b6-9671-9019b7b09bd0#e0e10aab-5977-49b3-a4dd-106db8821b6a
```

