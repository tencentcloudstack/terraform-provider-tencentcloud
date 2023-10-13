---
subcategory: "TencentCloud Managed Service for Grafana(TCMG)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_grafana_whitelist_config"
sidebar_current: "docs-tencentcloud-resource-monitor_grafana_whitelist_config"
description: |-
  Provides a resource to create a monitor grafana_whitelist_config
---

# tencentcloud_monitor_grafana_whitelist_config

Provides a resource to create a monitor grafana_whitelist_config

## Example Usage

```hcl
resource "tencentcloud_monitor_grafana_whitelist_config" "grafana_whitelist_config" {
  instance_id = "grafana-dp2hnnfa"
  whitelist   = ["10.1.1.1", "10.1.1.2", "10.1.1.3"]
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Grafana instance ID.
* `whitelist` - (Optional, Set: [`String`]) The addresses in the whitelist.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

monitor grafana_whitelist_config can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_grafana_whitelist_config.grafana_whitelist_config instance_id
```

