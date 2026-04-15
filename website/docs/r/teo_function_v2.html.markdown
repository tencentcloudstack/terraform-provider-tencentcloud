---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_function_v2"
sidebar_current: "docs-tencentcloud-resource-teo_function_v2"
description: |-
  Provides a resource to create a teo teo_function_v2
---

# tencentcloud_teo_function_v2

Provides a resource to create a teo teo_function_v2

## Example Usage

```hcl
resource "tencentcloud_teo_function_v2" "teo_function_v2" {
  content = <<-EOT
        addEventListener('fetch', e => {
          const response = new Response('Hello World!!');
          e.respondWith(response);
        });
    EOT
  name    = "aaa-zone-2qtuhspy7cr6-1310708577"
  remark  = "test"
  zone_id = "zone-2qtuhspy7cr6"
}
```

### Update Example

```hcl
resource "tencentcloud_teo_function_v2" "teo_function_v2" {
  content = <<-EOT
        addEventListener('fetch', e => {
          const response = new Response('Hello World');
          e.respondWith(response);
        });
    EOT
  name    = "aaa-zone-2qtuhspy7cr6-1310708577"
  remark  = "test-update"
  zone_id = "zone-2qtuhspy7cr6"
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required, String) Function content, currently only supports JavaScript code, with a maximum size of 5MB.
* `name` - (Required, String, ForceNew) Function name. It can only contain lowercase letters, numbers, hyphens, must start and end with a letter or number, and can have a maximum length of 30 characters.
* `zone_id` - (Required, String, ForceNew) ID of the site.
* `remark` - (Optional, String) Function description, maximum support of 60 characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time. The time is in Coordinated Universal Time (UTC) and follows the date and time format specified by the ISO 8601 standard.
* `domain` - The default domain name for the function.
* `function_id` - ID of the Function.
* `update_time` - Modification time. The time is in Coordinated Universal Time (UTC) and follows the date and time format specified by the ISO 8601 standard.


## Import

teo teo_function_v2 can be imported using id, e.g.

```
terraform import tencentcloud_teo_function_v2.teo_function_v2 zone_id#function_id
```

