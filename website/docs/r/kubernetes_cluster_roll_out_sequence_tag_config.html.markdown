---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_cluster_roll_out_sequence_tag_config"
sidebar_current: "docs-tencentcloud-resource-kubernetes_cluster_roll_out_sequence_tag_config"
description: |-
  Provides a resource to create a TKE kubernetes cluster roll out sequence tag config.
---

# tencentcloud_kubernetes_cluster_roll_out_sequence_tag_config

Provides a resource to create a TKE kubernetes cluster roll out sequence tag config.

## Example Usage

```hcl
resource "tencentcloud_kubernetes_cluster_roll_out_sequence_tag_config" "example" {
  cluster_id = "cls-6ml650mu"

  tags {
    key   = "Env"
    value = "Test"
  }

  tags {
    key   = "Protection-Level"
    value = "Medium"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) Cluster ID.
* `tags` - (Required, List) Cluster roll-out sequence tags. Supported tags: key `Env` with values [`Test`, `Pre-Production`, `Production`]; key `Protection-Level` with values [`Low`, `Medium`, `High`].

The `tags` object supports the following:

* `key` - (Required, String) Tag key.
* `value` - (Required, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

TKE kubernetes cluster roll out sequence tag config can be imported using the cluster_id, e.g.

```
terraform import tencentcloud_kubernetes_cluster_roll_out_sequence_tag_config.example cls-6ml650mu
```

