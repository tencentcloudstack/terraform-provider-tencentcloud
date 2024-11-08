---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_identity_center_scim_credential_status"
sidebar_current: "docs-tencentcloud-resource-identity_center_scim_credential_status"
description: |-
  Provides a resource to manage identity center scim credential status
---

# tencentcloud_identity_center_scim_credential_status

Provides a resource to manage identity center scim credential status

## Example Usage

```hcl
resource "tencentcloud_identity_center_scim_credential_status" "identity_center_scim_credential_status" {
  zone_id       = "z-xxxxxx"
  credential_id = "scimcred-xxxxxx"
  status        = "Enabled"
}
```

## Argument Reference

The following arguments are supported:

* `credential_id` - (Required, String, ForceNew) SCIM key ID. scimcred-prefix and followed by 12 random digits/lowercase letters.
* `status` - (Required, String) SCIM key status. Enabled-enabled. Disabled-disabled.
* `zone_id` - (Required, String, ForceNew) Space ID. z-prefix starts with 12 random digits/lowercase letters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

organization identity_center_scim_credential_status can be imported using the id, e.g.

```
terraform import tencentcloud_identity_center_scim_credential_status.identity_center_scim_credential_status ${zone_id}#${credential_id}
```

