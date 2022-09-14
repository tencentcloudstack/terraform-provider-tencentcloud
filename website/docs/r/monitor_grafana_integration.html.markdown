---
subcategory: "Monitor"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_grafana_integration"
sidebar_current: "docs-tencentcloud-resource-monitor_grafana_integration"
description: |-
  Provides a resource to create a monitor grafanaIntegration
---

# tencentcloud_monitor_grafana_integration

Provides a resource to create a monitor grafanaIntegration

## Example Usage



## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) grafana instance id.
* `content` - (Optional, String) generated json string of given integration json schema.
* `description` - (Optional, String) integration desc.
* `kind` - (Optional, String) integration json schema kind.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `integration_id` - integration id.


