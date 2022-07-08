---
subcategory: "Monitor"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_tmp_instance"
sidebar_current: "docs-tencentcloud-resource-monitor_tmp_instance"
description: |-
  Provides a resource to create a monitor tmpInstance
---

# tencentcloud_monitor_tmp_instance

Provides a resource to create a monitor tmpInstance

## Example Usage

```hcl
resource "tencentcloud_monitor_tmp_instance" "tmpInstance" {
  instance_name       = "logset-hello"
  vpc_id              = "vpc-2hfyray3"
  subnet_id           = "subnet-rdkj0agk"
  data_retention_time = 30
  zone                = "ap-guangzhou-3"
}
```

## Argument Reference

The following arguments are supported:

* `data_retention_time` - (Required) Data retention time.
* `instance_name` - (Required) Instance name.
* `subnet_id` - (Required) Subnet Id.
* `vpc_id` - (Required) Vpc Id.
* `zone` - (Required) Available zone.
* `grafana_instance_id` - (Optional) Associated grafana instance id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

monitor tmp instance can be imported using the id, e.g.
```
$ terraform import tencentcloud_monitor_tmp_instance.tmpInstance tmpInstance_id
```

