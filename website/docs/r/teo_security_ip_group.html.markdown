---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_security_ip_group"
sidebar_current: "docs-tencentcloud-resource-teo_security_ip_group"
description: |-
  Provides a resource to create a teo teo_security_ip_group
---

# tencentcloud_teo_security_ip_group

Provides a resource to create a teo teo_security_ip_group

## Example Usage

```hcl
resource "tencentcloud_teo_security_ip_group" "teo_security_ip_group" {
  zone_id = "zone-2qtuhspy7cr6"
  ip_group {
    content = [
      "10.1.1.1",
      "10.1.1.2",
      "10.1.1.3",
    ]
    name = "bbbbb"
  }
}
```

## Argument Reference

The following arguments are supported:

* `ip_group` - (Required, List) IP group information, replace all when modifying.
* `zone_id` - (Required, String) Site ID.

The `ip_group` object supports the following:

* `content` - (Required, Set) IP group content. Only supports IP and IP mask.
* `name` - (Required, String) Group name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

teo teo_security_ip_group can be imported using the id, e.g.

```
terraform import tencentcloud_teo_security_ip_group.teo_security_ip_group zone_id#group_id
```

