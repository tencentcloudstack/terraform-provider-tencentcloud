---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_ssl"
sidebar_current: "docs-tencentcloud-resource-redis_ssl"
description: |-
  Provides a resource to create a redis ssl
---

# tencentcloud_redis_ssl

Provides a resource to create a redis ssl

## Example Usage

```hcl
resource "tencentcloud_redis_ssl" "ssl" {
  instance_id = "crs-c1nl9rpv"
  ssl_config  = "disabled"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) The ID of instance.
* `ssl_config` - (Required, String) The SSL configuration status of the instance: `enabled`,`disabled`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

redis ssl can be imported using the instanceId, e.g.

```
terraform import tencentcloud_redis_ssl.ssl crs-c1nl9rpv
```

