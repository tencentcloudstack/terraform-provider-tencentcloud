---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_upgrade_multi_zone_operation"
sidebar_current: "docs-tencentcloud-resource-redis_upgrade_multi_zone_operation"
description: |-
  Provides a resource to create a redis upgrade_multi_zone_operation
---

# tencentcloud_redis_upgrade_multi_zone_operation

Provides a resource to create a redis upgrade_multi_zone_operation

## Example Usage

```hcl
resource "tencentcloud_redis_upgrade_multi_zone_operation" "upgrade_multi_zone_operation" {
  instance_id                    = "crs-c1nl9rpv"
  upgrade_proxy_and_redis_server = true
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) The ID of instance.
* `upgrade_proxy_and_redis_server` - (Optional, Bool, ForceNew) After you upgrade Multi-AZ, whether the nearby access feature is supported.true: Supports nearby access.The upgrade process, which requires upgrading both the proxy version and the Redis kernel minor version, involves data migration and can take several hours.false: No need to support nearby access.Upgrading Multi-AZ only involves managing metadata migration, with no service impact, and the upgrade process typically completes within 3 minutes.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



