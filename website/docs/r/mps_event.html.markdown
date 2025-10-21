---
subcategory: "Media Processing Service(MPS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mps_event"
sidebar_current: "docs-tencentcloud-resource-mps_event"
description: |-
  Provides a resource to create a mps event
---

# tencentcloud_mps_event

Provides a resource to create a mps event

## Example Usage

```hcl
resource "tencentcloud_mps_event" "event" {
  event_name  = "you-event-name"
  description = "event description"
}
```

## Argument Reference

The following arguments are supported:

* `event_name` - (Required, String) Event name.
* `description` - (Optional, String) Event description.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mps event can be imported using the id, e.g.

```
terraform import tencentcloud_mps_event.event event_id
```

