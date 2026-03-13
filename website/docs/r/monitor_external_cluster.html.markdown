---
subcategory: "Cloud Monitor(Monitor)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_external_cluster"
sidebar_current: "docs-tencentcloud-resource-monitor_external_cluster"
description: |-
  Provides a resource to create a Monitor external cluster
---

# tencentcloud_monitor_external_cluster

Provides a resource to create a Monitor external cluster

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
```

## Argument Reference

The following arguments are supported:

* `cluster_region` - (Required, String, ForceNew) Cluster region.
* `instance_id` - (Required, String, ForceNew) Instance ID.
* `cluster_name` - (Optional, String, ForceNew) Cluster name.
* `enable_external` - (Optional, Bool, ForceNew) Whether to enable public network access.
* `external_labels` - (Optional, List, ForceNew) External labels.

The `external_labels` object supports the following:

* `name` - (Required, String) Label name.
* `value` - (Optional, String) Label value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cluster_id` - Cluster ID.
* `cluster_type` - Cluster type, returned by API.
* `status` - Agent status.


## Import

Monitor external cluster can be imported using the instanceId#clusterId, e.g.

```
terraform import tencentcloud_monitor_external_cluster.example prom-gzg3f1em#ecls-qi9v5opk
```

