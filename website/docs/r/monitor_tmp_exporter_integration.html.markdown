---
subcategory: "Monitor"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_tmp_exporter_integration"
sidebar_current: "docs-tencentcloud-resource-monitor_tmp_exporter_integration"
description: |-
  Provides a resource to create a monitor tmpExporterIntegration
---

# tencentcloud_monitor_tmp_exporter_integration

Provides a resource to create a monitor tmpExporterIntegration

## Example Usage

```hcl
resource "tencentcloud_monitor_tmp_exporter_integration" "basic" {
  instance_id = "prom-dko9d0nu"
  kind        = "cvm-http-sd-exporter"
  content     = "{\"kind\":\"cvm-http-sd-exporter\",\"spec\":{\"job\":\"job_name: example-job-name-test\\nmetrics_path: /metrics\\ncvm_sd_configs:\\n- region: ap-guangzhou\\n  ports:\\n  - 9100\\n  filters:         \\n  - name: tag:示例标签键\\n    values: \\n    - 示例标签值\\nrelabel_configs: \\n- source_labels: [__meta_cvm_instance_state]\\n  regex: RUNNING\\n  action: keep\\n- regex: __meta_cvm_tag_(.*)\\n  replacement: $1\\n  action: labelmap\\n- source_labels: [__meta_cvm_region]\\n  target_label: region\\n  action: replace\"}}"
  kube_type   = 1
  cluster_id  = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) Cluster ID.
* `content` - (Required, String) Integration config.
* `instance_id` - (Required, String) Instance id.
* `kind` - (Required, String) Type.
* `kube_type` - (Required, Int) Integration config.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



