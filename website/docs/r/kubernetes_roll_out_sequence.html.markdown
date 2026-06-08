---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_roll_out_sequence"
sidebar_current: "docs-tencentcloud-resource-kubernetes_roll_out_sequence"
description: |-
  Provides a resource to create a TKE roll-out sequence
---

# tencentcloud_kubernetes_roll_out_sequence

Provides a resource to create a TKE roll-out sequence

## Example Usage

```hcl
resource "tencentcloud_kubernetes_roll_out_sequence" "example" {
  name    = "tf-example"
  enabled = true

  sequence_flows {
    tags {
      key   = "env"
      value = ["production"]
    }

    soak_time = 3600
  }

  sequence_flows {
    tags {
      key   = "env"
      value = ["staging", "testing"]
    }

    soak_time = 1800
  }
}
```

## Argument Reference

The following arguments are supported:

* `enabled` - (Required, Bool) Whether the roll-out sequence is enabled.
* `name` - (Required, String) The name of the roll-out sequence.
* `sequence_flows` - (Required, List) The sequence flow steps of the roll-out sequence.

The `sequence_flows` object supports the following:

* `soak_time` - (Required, Int) Wait time in seconds between steps.
* `tags` - (Required, List) The tags for the sequence flow step.

The `tags` object of `sequence_flows` supports the following:

* `key` - (Required, String) Tag key.
* `value` - (Required, List) Tag values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

TKE roll-out sequence can be imported using the id, e.g.

```
terraform import tencentcloud_kubernetes_roll_out_sequence.example 123
```

