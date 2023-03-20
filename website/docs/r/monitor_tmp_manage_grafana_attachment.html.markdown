---
subcategory: "Managed Service for Prometheus(TMP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_tmp_manage_grafana_attachment"
sidebar_current: "docs-tencentcloud-resource-monitor_tmp_manage_grafana_attachment"
description: |-
  Provides a resource to create a monitor tmp_manage_grafana_attachment
---

# tencentcloud_monitor_tmp_manage_grafana_attachment

Provides a resource to create a monitor tmp_manage_grafana_attachment

## Example Usage

```hcl
resource "tencentcloud_monitor_tmp_manage_grafana_attachment" "manage_grafana_attachment" {
  grafana_id  = "grafana-xxxxxx"
  instance_id = "prom-xxxxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `grafana_id` - (Required, String, ForceNew) Grafana instance ID.
* `instance_id` - (Required, String, ForceNew) Prometheus instance ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

monitor tmp_manage_grafana_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_tmp_manage_grafana_attachment.manage_grafana_attachment prom-xxxxxxxx
```

