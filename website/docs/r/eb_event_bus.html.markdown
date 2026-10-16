---
subcategory: "EventBridge(EB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eb_event_bus"
sidebar_current: "docs-tencentcloud-resource-eb_event_bus"
description: |-
  Provides a resource to create a eb event_bus
---

# tencentcloud_eb_event_bus

Provides a resource to create a eb event_bus

## Example Usage

```hcl
resource "tencentcloud_eb_event_bus" "foo" {
  event_bus_name = "tf-event_bus"
  description    = "event bus desc"
  enable_store   = false
  save_days      = 1
  tags = {
    "createdBy" = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `event_bus_name` - (Required, String) Event set name, which can only contain letters, numbers, underscores, hyphens, starts with a letter and ends with a number or letter, 2~60 characters.
* `description` - (Optional, String) Event set description, unlimited character type, description within 200 characters.
* `enable_store` - (Optional, Bool) Whether the EB storage is enabled.
* `save_days` - (Optional, Int) EB storage duration.
* `tags` - (Optional, Map) Tag description list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateEventBus](https://www.tencentcloud.com/document/api/1108/45829)
* Create: [ModifyResourceTags](https://www.tencentcloud.com/document/api/651/33012)
* Read: [DescribeResourceTagsByResourceIds](https://www.tencentcloud.com/document/api/651/33022)
* Read: [GetEventBus](https://www.tencentcloud.com/document/api/1108/45827)
* Update: [GetEventBus](https://www.tencentcloud.com/document/api/1108/45827)
* Update: [ModifyResourceTags](https://www.tencentcloud.com/document/api/651/33012)
* Update: [UpdateEventBus](https://www.tencentcloud.com/document/api/1108/45825)
* Delete: [DeleteEventBus](https://www.tencentcloud.com/document/api/1108/45828)

## Import

eb event_bus can be imported using the id, e.g.

```
terraform import tencentcloud_eb_event_bus.event_bus event_bus_id
```

