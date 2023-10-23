---
subcategory: "TencentCloud Managed Service for Grafana(TCMG)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_grafana_dns_config"
sidebar_current: "docs-tencentcloud-resource-monitor_grafana_dns_config"
description: |-
  Provides a resource to create a monitor grafana_dns_config
---

# tencentcloud_monitor_grafana_dns_config

Provides a resource to create a monitor grafana_dns_config

## Example Usage

```hcl
resource "tencentcloud_monitor_grafana_dns_config" "grafana_dns_config" {
  instance_id  = "grafana-dp2hnnfa"
  name_servers = ["10.1.2.1", "10.1.2.2", "10.1.2.3"]
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Grafana instance ID.
* `name_servers` - (Optional, Set: [`String`]) DNS nameserver list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

monitor grafana_dns_config can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_grafana_dns_config.grafana_dns_config instance_id
```

