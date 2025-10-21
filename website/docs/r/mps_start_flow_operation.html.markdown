---
subcategory: "Media Processing Service(MPS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mps_start_flow_operation"
sidebar_current: "docs-tencentcloud-resource-mps_start_flow_operation"
description: |-
  Provides a resource to create a mps start_flow_operation
---

# tencentcloud_mps_start_flow_operation

Provides a resource to create a mps start_flow_operation

## Example Usage

### Start flow

```hcl
resource "tencentcloud_mps_start_flow_operation" "operation" {
  flow_id = tencentcloud_mps_flow.flow_rtp.id
  start   = true
}
```

### Stop flow

```hcl
resource "tencentcloud_mps_start_flow_operation" "operation" {
  flow_id = tencentcloud_mps_flow.flow_rtp.id
  start   = false
}
```

## Argument Reference

The following arguments are supported:

* `flow_id` - (Required, String, ForceNew) Flow Id.
* `start` - (Required, Bool, ForceNew) `true`: start mps stream link flow; `false`: stop.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



