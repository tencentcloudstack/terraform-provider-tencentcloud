---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_log_delivery"
sidebar_current: "docs-tencentcloud-resource-redis_log_delivery"
description: |-
  Provides a resource to create Redis instance log delivery land set its attributes.
---

# tencentcloud_redis_log_delivery

Provides a resource to create Redis instance log delivery land set its attributes.

~> **NOTE:** When you use an existing cls logset and topic to enable logging, there is no need to set parameters such
as `period`, `create_index`, `log_region`, etc.

## Example Usage

### Use cls logset and topic which existed

```hcl
resource "tencentcloud_redis_log_delivery" "delivery" {
  instance_id = "crs-dmjj8en7"
  logset_id   = "cc31d9d6-74c0-4888-8b2f-b8148c3bcc5c"
  topic_id    = "5c2333e9-0bab-41fd-9f75-c602b3f9545f"
}
```

### Use exist cls logset and create new topic

```hcl
resource "tencentcloud_redis_log_delivery" "delivery" {
  instance_id  = "crs-dmjj8en7"
  logset_id    = "cc31d9d6-74c0-4888-8b2f-b8148c3bcc5c"
  topic_name   = "test13"
  period       = 20
  create_index = true
}
```

### Create new cls logset and topic

```hcl
resource "tencentcloud_redis_log_delivery" "delivery" {
  instance_id  = "crs-dmjj8en7"
  log_region   = "ap-guangzhou"
  logset_name  = "test"
  topic_name   = "test"
  period       = 20
  create_index = true
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID.
* `create_index` - (Optional, Bool) Whether to create an index when creating a log topic.
* `log_region` - (Optional, String) The region where the log set is located; if not specified, the region where the instance is located will be used by default.
* `logset_id` - (Optional, String) The ID of the log set being delivered.
* `logset_name` - (Optional, String) Log set name. If LogsetId does not specify a specific log set ID, please configure this parameter to set the log set name, and the system will automatically create a new log set with the specified name.
* `period` - (Optional, Int) Log storage time, defaults to 30 days, with an optional range of 1-3600 days.
* `topic_id` - (Optional, String) The ID of the topic being delivered.
* `topic_name` - (Optional, String) Log topic name, required when TopicId is empty, a new log topic will be automatically created.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Redis log delivery can be imported, e.g.

```
$ terraform import tencentcloud_redis_log_delivery.delivery crs-dmjj8en7
```

