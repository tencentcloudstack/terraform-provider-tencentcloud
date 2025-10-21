---
subcategory: "TencentCloud Managed Service for Grafana(TCMG)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_grafana_version_upgrade"
sidebar_current: "docs-tencentcloud-resource-monitor_grafana_version_upgrade"
description: |-
  Provides a resource to create a monitor grafana_version_upgrade
---

# tencentcloud_monitor_grafana_version_upgrade

Provides a resource to create a monitor grafana_version_upgrade

## Example Usage

```hcl
resource "tencentcloud_monitor_grafana_version_upgrade" "grafana_version_upgrade" {
  instance_id = "grafana-dp2hnnfa"
  alias       = "v8.2.7"
}
```

## Argument Reference

The following arguments are supported:

* `alias` - (Required, String) Version alias.
* `instance_id` - (Required, String, ForceNew) Grafana instance ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

monitor grafana_version_upgrade can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_grafana_version_upgrade.grafana_version_upgrade instance_id
```

