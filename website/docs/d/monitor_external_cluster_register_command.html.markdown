---
subcategory: "Cloud Monitor(Monitor)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_external_cluster_register_command"
sidebar_current: "docs-tencentcloud-datasource-monitor_external_cluster_register_command"
description: |-
  Use this data source to query Monitor external cluster register command
---

# tencentcloud_monitor_external_cluster_register_command

Use this data source to query Monitor external cluster register command

## Example Usage

```hcl
resource "tencentcloud_monitor_external_cluster" "example" {
  instance_id    = "prom-gzg3f1em"
  cluster_region = "ap-guangzhou"
  cluster_name   = "tf-external-cluster"

  external_labels {
    name  = "clusterName"
    value = "example"
  }

  external_labels {
    name  = "environment"
    value = "prod"
  }

  enable_external = false
}

data "tencentcloud_monitor_external_cluster_register_command" "example" {
  instance_id = tencentcloud_monitor_external_cluster.example.instance_id
  cluster_id  = tencentcloud_monitor_external_cluster.example.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) External cluster ID.
* `instance_id` - (Required, String) TMP instance ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `command` - Register command, contains YAML configuration for Kubernetes cluster installation.


