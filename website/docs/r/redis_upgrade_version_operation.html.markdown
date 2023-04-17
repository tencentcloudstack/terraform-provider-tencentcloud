---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_upgrade_version_operation"
sidebar_current: "docs-tencentcloud-resource-redis_upgrade_version_operation"
description: |-
  Provides a resource to create a redis upgrade_version_operation
---

# tencentcloud_redis_upgrade_version_operation

Provides a resource to create a redis upgrade_version_operation

## Example Usage

```hcl
resource "tencentcloud_redis_upgrade_version_operation" "upgrade_version_operation" {
  instance_id          = "crs-c1nl9rpv"
  target_instance_type = "6"
  switch_option        = 2
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) The ID of instance.
* `switch_option` - (Required, Int, ForceNew) Switch mode:1 - maintenance time window switching,2 - immediate switching.
* `target_instance_type` - (Required, String, ForceNew) Target instance type, same as [CreateInstances](https://cloud.tencent.com/document/api/239/20026), that is, the target type of the instance to be changed.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



