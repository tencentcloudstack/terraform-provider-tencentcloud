---
subcategory: "Monitor"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_grafana_instance"
sidebar_current: "docs-tencentcloud-resource-monitor_grafana_instance"
description: |-
  Provides a resource to create a monitor grafanaInstance
---

# tencentcloud_monitor_grafana_instance

Provides a resource to create a monitor grafanaInstance

## Example Usage

```hcl
resource "tencentcloud_monitor_grafana_instance" "grafanaInstance" {
  instance_name         = "test-grafana"
  vpc_id                = "vpc-2hfyray3"
  subnet_ids            = ["subnet-rdkj0agk"]
  grafana_init_password = "1234567890"
  enable_internet       = false

  tags = {
    "createdBy" = "test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_name` - (Required, String) Instance name.
* `enable_internet` - (Optional, Bool) Control whether grafana could be accessed by internet.
* `grafana_init_password` - (Optional, String) Grafana server admin password.
* `subnet_ids` - (Optional, Set: [`String`]) Subnet Id array.
* `tags` - (Optional, Map) Tag description list.
* `vpc_id` - (Optional, String) Vpc Id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `instance_id` - Grafana instance id.
* `instance_status` - Grafana instance status, 1: Creating, 2: Running, 6: Stopped.
* `root_url` - Grafana external url which could be accessed by user.


## Import

monitor grafanaInstance can be imported using the id, e.g.
```
$ terraform import tencentcloud_monitor_grafana_instance.grafanaInstance grafanaInstance_id
```

