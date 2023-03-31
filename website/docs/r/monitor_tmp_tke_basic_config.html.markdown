---
subcategory: "Cloud Monitor(Monitor)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_tmp_tke_basic_config"
sidebar_current: "docs-tencentcloud-resource-monitor_tmp_tke_basic_config"
description: |-
  Provides a resource to create a monitor tmp_tke_basic_config
---

# tencentcloud_monitor_tmp_tke_basic_config

Provides a resource to create a monitor tmp_tke_basic_config

## Example Usage

```hcl
resource "tencentcloud_monitor_tmp_tke_basic_config" "tmp_tke_basic_config" {
  instance_id  = "prom-xxxxxx"
  cluster_type = "eks"
  cluster_id   = "cls-xxxxxx"
  name         = "kube-system/kube-state-metrics"
  metrics_name = ["kube_job_status_succeeded"]
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) ID of cluster.
* `cluster_type` - (Required, String) Type of cluster.
* `instance_id` - (Required, String) ID of instance.
* `metrics_name` - (Required, Set: [`String`]) Configure the name of the metric to keep on.
* `name` - (Required, String) Name. The naming rule is: namespace/name. If you don&#39;t have any namespace, use the default namespace: kube-system, otherwise use the specified one.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `config_type` - config type, `service_monitors`, `pod_monitors`, `raw_jobs`.
* `config` - Full configuration in yaml format.


