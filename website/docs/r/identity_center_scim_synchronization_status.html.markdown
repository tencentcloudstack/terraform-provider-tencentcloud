---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_identity_center_scim_synchronization_status"
sidebar_current: "docs-tencentcloud-resource-identity_center_scim_synchronization_status"
description: |-
  Provides a resource to manage identity center scim synchronization status
---

# tencentcloud_identity_center_scim_synchronization_status

Provides a resource to manage identity center scim synchronization status

## Example Usage

```hcl
resource "tencentcloud_identity_center_scim_synchronization_status" "identity_center_scim_synchronization_status" {
  zone_id                     = "z-xxxxxx"
  scim_synchronization_status = "Enabled"
}
```

## Argument Reference

The following arguments are supported:

* `scim_synchronization_status` - (Required, String) SCIM synchronization status. Enabled-enabled. Disabled-disables.
* `zone_id` - (Required, String, ForceNew) Space ID. z-prefix starts with 12 random digits/lowercase letters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

organization identity_center_scim_synchronization_status can be imported using the id, e.g.

```
terraform import tencentcloud_identity_center_scim_synchronization_status.identity_center_scim_synchronization_status ${zone_id}
```

