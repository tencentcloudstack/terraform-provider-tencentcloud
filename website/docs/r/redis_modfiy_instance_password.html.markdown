---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_modfiy_instance_password"
sidebar_current: "docs-tencentcloud-resource-redis_modfiy_instance_password"
description: |-
  Provides a resource to create a redis modfiy_instance_password
---

# tencentcloud_redis_modfiy_instance_password

Provides a resource to create a redis modfiy_instance_password

## Example Usage

```hcl
resource "tencentcloud_redis_modfiy_instance_password" "modfiy_instance_password" {
  instance_id  = "crs-c1nl9rpv"
  old_password = ""
  password     = ""
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) The ID of instance.
* `old_password` - (Required, String) The old password for the instance.
* `password` - (Required, String) The password for the instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



