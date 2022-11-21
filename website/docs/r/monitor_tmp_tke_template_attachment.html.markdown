---
subcategory: "Cloud Monitor(Monitor)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_tmp_tke_template_attachment"
sidebar_current: "docs-tencentcloud-resource-monitor_tmp_tke_template_attachment"
description: |-
  Provides a resource to create a tmp tke template attachment
---

# tencentcloud_monitor_tmp_tke_template_attachment

Provides a resource to create a tmp tke template attachment

## Example Usage

```hcl
resource "tencentcloud_monitor_tmp_tke_template_attachment" "temp_attachment" {
  template_id = "temp-xxx"

  targets {
    region      = "ap-xxx"
    instance_id = "prom-xxx"
  }
}
```

## Argument Reference

The following arguments are supported:

* `targets` - (Required, List, ForceNew) Sync target details.
* `template_id` - (Required, String, ForceNew) The ID of the template, which is used for the outgoing reference.

The `targets` object supports the following:

* `instance_id` - (Required, String) instance id.
* `region` - (Required, String) target area.
* `cluster_id` - (Optional, String) ID of the cluster.
* `cluster_name` - (Optional, String) Name the cluster.
* `cluster_type` - (Optional, String) Cluster type.
* `instance_name` - (Optional, String) Name of the prometheus instance.
* `sync_time` - (Optional, String) Last sync template time.
* `version` - (Optional, String) Template version currently in use.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



