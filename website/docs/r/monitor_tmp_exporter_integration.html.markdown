---
subcategory: "Managed Service for Prometheus(TMP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_tmp_exporter_integration"
sidebar_current: "docs-tencentcloud-resource-monitor_tmp_exporter_integration"
description: |-
  Provides a resource to create a monitor tmpExporterIntegration
---

# tencentcloud_monitor_tmp_exporter_integration

Provides a resource to create a monitor tmpExporterIntegration

~> **NOTE:** If you only want to upgrade the exporter version with same config, you can set `version` under `instanceSpec` with any value to trigger the change.

## Example Usage

Use blackbox-exporter

```hcl
resource "tencentcloud_monitor_tmp_exporter_integration" "tmpExporterIntegration" {
  instance_id = "prom-dko9d0nu"
  kind        = "blackbox-exporter"
  content     = "{\"name\":\"test\",\"kind\":\"blackbox-exporter\",\"spec\":{\"instanceSpec\":{\"module\":\"http_get\",\"urls\":[\"xx\"]}}}"
  kube_type   = 1
  cluster_id  = "cls-bmuaukfu"
}
```

Use es-exporter

```hcl
resource "tencentcloud_monitor_tmp_exporter_integration" "tmpExporterIntegrationEs" {
  instance_id = tencentcloud_monitor_tmp_instance.tmpInstance.id
  kind        = "es-exporter"
  content = jsonencode({
    "name" : "ex-exporter-example",
    "kind" : "es-exporter",
    "spec" : {
      "instanceSpec" : {
        "url" : "http://127.0.0.1:9123",
        "labels" : {
          "instance" : "es-abcd"
        },
        "version" : "1.70.1",
        "user" : "fugiat Duis minim",
        "password" : "exercitation cillum velit"
      },
      "exporterSpec" : {
        "all" : true,
        "indicesSettings" : false,
        "snapshots" : false,
        "indices" : true,
        "shards" : false
      }
    }
  })
  cluster_id = ""
  kube_type  = 3
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



