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

```hcl
resource "tencentcloud_monitor_tmp_tke_config" "basic" {
  instance_id  = "prom-1lspn8sw"
  cluster_type = "tke"
  cluster_id   = "cls-ely08ic4"
  service_monitors {
    name   = "service-monitor-001"
    config = "apiVersion: monitoring.coreos.com/v1\nkind: ServiceMonitor\nmetadata:\n  name: service-monitor-001\n  namespace: kube-system\nspec:\n  endpoints:\n    - interval: 115s\n      port: 8080-8080-tcp\n      path: /metrics\n      relabelings:\n        - action: replace\n          sourceLabels:\n            - __meta_kubernetes_pod_label_app\n          targetLabel: application\n  namespaceSelector:\n    matchNames:\n      - test\n  selector:\n    matchLabels:\n      app: test"
  }

  pod_monitors {
    name   = "pod-monitor-001"
    config = "apiVersion: monitoring.coreos.com/v1\nkind: PodMonitor\nmetadata:\n  name: pod-monitor-001\n  namespace: kube-system\nspec:\n  podMetricsEndpoints:\n    - interval: 15s\n      port: metric-port\n      path: /metrics\n      relabelings:\n        - action: replace\n          sourceLabels:\n            - instance\n          regex: (.*)\n          targetLabel: instance\n          replacement: xxxxxx\n  namespaceSelector:\n    matchNames:\n      - test\n  selector:\n    matchLabels:\n      k8s-app: test"
  }

  raw_jobs {
    name   = "raw_jobs_001"
    config = "scrape_configs:\n- job_name: raw_jobs_001\n  honor_labels: true\n"
  }
}
```

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


