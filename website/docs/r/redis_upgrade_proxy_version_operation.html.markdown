---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_upgrade_proxy_version_operation"
sidebar_current: "docs-tencentcloud-resource-redis_upgrade_proxy_version_operation"
description: |-
  Provides a resource to create a redis upgrade proxy version operation
---

# tencentcloud_redis_upgrade_proxy_version_operation

Provides a resource to create a redis upgrade proxy version operation

## Example Usage

```hcl
resource "tencentcloud_redis_upgrade_proxy_version_operation" "example" {
  instance_id               = "crs-c1nl9rpv"
  current_proxy_version     = "5.0.0"
  upgrade_proxy_version     = "5.8.12"
  instance_type_upgrade_now = 1
}
```

## Argument Reference

The following arguments are supported:

* `current_proxy_version` - (Required, String, ForceNew) Current proxy version.
* `instance_id` - (Required, String, ForceNew) The ID of instance.
* `instance_type_upgrade_now` - (Required, Int, ForceNew) Switch mode:1 - Upgrade now0 - Maintenance window upgrade.
* `upgrade_proxy_version` - (Required, String, ForceNew) Upgradeable redis proxy version.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



