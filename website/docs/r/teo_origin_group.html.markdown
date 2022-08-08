---
subcategory: "Teo"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_origin_group"
sidebar_current: "docs-tencentcloud-resource-teo_origin_group"
description: |-
  Provides a resource to create a teo originGroup
---

# tencentcloud_teo_origin_group

Provides a resource to create a teo originGroup

## Example Usage

```hcl
resource "tencentcloud_teo_origin_group" "originGroup" {
  record {
    private_parameter {}
  }
  tags = {
    "createdBy" = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `origin_id` - (Required, String) .
* `origin_name` - (Required, String) .
* `record` - (Required, List) .
* `type` - (Required, String) area, weight.
* `zone_id` - (Required, String) .
* `origin_type` - (Optional, String) .
* `tags` - (Optional, Map) Tag description list.

The `private_parameter` object supports the following:

* `name` - (Required, String) .
* `value` - (Required, String) .

The `record` object supports the following:

* `area` - (Required, Set) .
* `port` - (Required, Int) .
* `record` - (Required, String) .
* `weight` - (Required, Int) 1-100.
* `private_parameter` - (Optional, List) .
* `private` - (Optional, Bool) .

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `application_proxy_used` - .
* `load_balancing_used_type` - none, dns_only, proxied, both.
* `load_balancing_used` - .
* `update_time` - .
* `zone_name` - .


## Import

teo originGroup can be imported using the id, e.g.
```
$ terraform import tencentcloud_teo_origin_group.originGroup originGroup_id
```

