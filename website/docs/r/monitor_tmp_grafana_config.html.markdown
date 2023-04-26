---
subcategory: "TencentCloud Managed Service for Grafana(TCMG)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_tmp_grafana_config"
sidebar_current: "docs-tencentcloud-resource-monitor_tmp_grafana_config"
description: |-
  Provides a resource to create a monitor tmp_grafana_config
---

# tencentcloud_monitor_tmp_grafana_config

Provides a resource to create a monitor tmp_grafana_config

## Example Usage

```hcl
resource "tencentcloud_monitor_tmp_grafana_config" "tmp_grafana_config" {
  config = jsonencode(
    {
      server = {
        http_port           = 8080
        root_url            = "https://cloud-grafana.woa.com/grafana-ffrdnrfa/"
        serve_from_sub_path = true
      }
    }
  )
  instance_id = "grafana-29phe08q"
}
```

## Argument Reference

The following arguments are supported:

* `config` - (Required, String) JSON encoded string.
* `instance_id` - (Required, String) Instance id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

monitor tmp_grafana_config can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_tmp_grafana_config.tmp_grafana_config tmp_grafana_config_id
```

