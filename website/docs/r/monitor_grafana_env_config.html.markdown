---
subcategory: "TencentCloud Managed Service for Grafana(TCMG)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_grafana_env_config"
sidebar_current: "docs-tencentcloud-resource-monitor_grafana_env_config"
description: |-
  Provides a resource to create a monitor grafana_env_config
---

# tencentcloud_monitor_grafana_env_config

Provides a resource to create a monitor grafana_env_config

## Example Usage

```hcl
resource "tencentcloud_monitor_grafana_env_config" "grafana_env_config" {
  instance_id = "grafana-dp2hnnfa"
  envs = {
    "aaa" = "ccc"
    "bbb" = "ccc"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Grafana instance ID.
* `envs` - (Optional, Map) Environment variables.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

monitor grafana_env_config can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_grafana_env_config.grafana_env_config instance_id
```

