---
subcategory: "Media Processing Service(MPS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mps_parse_live_stream_process_notify_operation"
sidebar_current: "docs-tencentcloud-resource-mps_parse_live_stream_process_notify_operation"
description: |-
  Provides a resource to create a mps parse_live_stream_process_notify_operation
---

# tencentcloud_mps_parse_live_stream_process_notify_operation

Provides a resource to create a mps parse_live_stream_process_notify_operation

## Example Usage

```hcl
resource "tencentcloud_mps_parse_live_stream_process_notify_operation" "operation" {
  content = "{\"EventType\":\"WorkflowTask\", xxx}"
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required, String, ForceNew) Live stream event notification obtained from CMQ.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



