---
subcategory: "DNSPOD"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dnspod_snapshot_config"
sidebar_current: "docs-tencentcloud-resource-dnspod_snapshot_config"
description: |-
  Provides a resource to create a dnspod snapshot_config
---

# tencentcloud_dnspod_snapshot_config

Provides a resource to create a dnspod snapshot_config

## Example Usage

```hcl
resource "tencentcloud_dnspod_snapshot_config" "snapshot_config" {
  domain = "dnspod.cn"
  period = "hourly"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) Domain name.
* `period` - (Required, String) Backup interval: empty string - no backup, half_hour - every half hour, hourly - every hour, daily - every day, monthly - every month.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

dnspod snapshot_config can be imported using the id, e.g.

```
terraform import tencentcloud_dnspod_snapshot_config.snapshot_config domain
```

