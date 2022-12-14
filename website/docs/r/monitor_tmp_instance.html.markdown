---
subcategory: "Managed Service for Prometheus(TMP)"
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
  instance_name       = "demo"
  vpc_id              = "vpc-2hfyray3"
  subnet_id           = "subnet-rdkj0agk"
  data_retention_time = 30
  zone                = "ap-guangzhou-3"
  tags = {
    "createdBy" = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `data_retention_time` - (Required, Int) Data retention time.
* `instance_name` - (Required, String) Instance name.
* `subnet_id` - (Required, String) Subnet Id.
* `vpc_id` - (Required, String) Vpc Id.
* `zone` - (Required, String) Available zone.
* `tags` - (Optional, Map) Tag description list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

monitor tmpInstance can be imported using the id, e.g.
```
$ terraform import tencentcloud_monitor_tmp_instance.tmpInstance tmpInstance_id
```

