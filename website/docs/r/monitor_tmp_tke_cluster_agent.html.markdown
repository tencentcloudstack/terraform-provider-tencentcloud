---
subcategory: "Monitor"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_tmp_tke_cluster_agent"
sidebar_current: "docs-tencentcloud-resource-monitor_tmp_tke_cluster_agent"
description: |-
  Provides a resource to create a tmp tke cluster agent
---

# tencentcloud_monitor_tmp_tke_cluster_agent

Provides a resource to create a tmp tke cluster agent

## Example Usage

```hcl
resource "tencentcloud_monitor_tmp_tke_cluster_agent" "tmpClusterAgent" {
  instance_id = "prom-xxx"

  agents {
    region          = "ap-xxx"
    cluster_type    = "eks"
    cluster_id      = "cls-xxx"
    enable_external = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `agents` - (Required, List) agent list.
* `instance_id` - (Required, String, ForceNew) Instance Id.

The `agents` object supports the following:

* `cluster_id` - (Required, String) An id identify the cluster, like `cls-xxxxxx`.
* `cluster_type` - (Required, String) Type of cluster.
* `enable_external` - (Required, Bool) Whether to enable the public network CLB.
* `region` - (Required, String) Limitation of region.
* `external_labels` - (Optional, List) All metrics collected by the cluster will carry these labels.
* `in_cluster_pod_config` - (Optional, Map) Pod configuration for components deployed in the cluster.
* `not_install_basic_scrape` - (Optional, Bool) Whether to install the default collection configuration.
* `not_scrape` - (Optional, Bool) Whether to collect indicators, true means drop all indicators, false means collect default indicators.

The `external_labels` object supports the following:

* `name` - (Required, String) Indicator name.
* `value` - (Optional, String) Index value.

The `in_cluster_pod_config` object supports the following:

* `host_net` - (Required, Bool) Whether to use HostNetWork.
* `node_selector` - (Optional, List) Specify the pod to run the node.
* `tolerations` - (Optional, List) Tolerate Stain.

The `node_selector` object supports the following:

* `name` - (Optional, String) The pod configuration name of the component deployed in the cluster.
* `value` - (Optional, String) Pod configuration values for components deployed in the cluster.

The `tolerations` object supports the following:

* `effect` - (Optional, String) blemish effect to match.
* `key` - (Optional, String) The taint key to which the tolerance applies.
* `operator` - (Optional, String) key-value relationship.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



