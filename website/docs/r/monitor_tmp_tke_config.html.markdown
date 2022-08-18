---
subcategory: "Monitor"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_tmp_tke_config"
sidebar_current: "docs-tencentcloud-resource-monitor_tmp_tke_config"
description: |-
  Provides a resource to create a tke tmpPrometheusConfig
---

# tencentcloud_monitor_tmp_tke_config

Provides a resource to create a tke tmpPrometheusConfig

## Example Usage



## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) ID of cluster.
* `cluster_type` - (Required, String, ForceNew) Type of cluster.
* `instance_id` - (Required, String, ForceNew) ID of instance.
* `pod_monitors` - (Optional, List) Configuration of the pod monitors.
* `raw_jobs` - (Optional, List) Configuration of the native prometheus job.
* `service_monitors` - (Optional, List) Configuration of the service monitors.

The `pod_monitors` object supports the following:

* `config` - (Required, String) Config.
* `name` - (Required, String) Name. The naming rule is: namespace/name. If you don't have any namespace, use the default namespace: kube-system, otherwise use the specified one.
* `template_id` - (Optional, String) Used for output parameters, if the configuration comes from a template, it is the template id.

The `raw_jobs` object supports the following:

* `config` - (Required, String) Config.
* `name` - (Required, String) Name.
* `template_id` - (Optional, String) Used for output parameters, if the configuration comes from a template, it is the template id.

The `service_monitors` object supports the following:

* `config` - (Required, String) Config.
* `name` - (Required, String) Name. The naming rule is: namespace/name. If you don't have any namespace, use the default namespace: kube-system, otherwise use the specified one.
* `template_id` - (Optional, String) Used for output parameters, if the configuration comes from a template, it is the template id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `config` - Global configuration.


