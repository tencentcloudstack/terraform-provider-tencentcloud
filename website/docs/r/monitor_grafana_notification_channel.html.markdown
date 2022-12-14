---
subcategory: "TencentCloud Managed Service for Grafana(TCMG)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_grafana_notification_channel"
sidebar_current: "docs-tencentcloud-resource-monitor_grafana_notification_channel"
description: |-
  Provides a resource to create a monitor grafanaNotificationChannel
---

# tencentcloud_monitor_grafana_notification_channel

Provides a resource to create a monitor grafanaNotificationChannel

## Example Usage

```hcl
resource "tencentcloud_monitor_grafana_notification_channel" "grafanaNotificationChannel" {
  instance_id   = "grafana-50nj6v00"
  channel_name  = "create-channel"
  org_id        = 1
  receivers     = ["Consumer-6vkna7pevq"]
  extra_org_ids = []
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) grafana instance id.
* `channel_name` - (Optional, String) channel name.
* `extra_org_ids` - (Optional, Set: [`String`]) extra grafana organization id list, default to 1 representing Main Org.
* `org_id` - (Optional, Int) Grafana organization which channel will be installed, default to 1 representing Main Org.
* `receivers` - (Optional, Set: [`String`]) cloud monitor notification template notice-id list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `channel_id` - plugin id.


