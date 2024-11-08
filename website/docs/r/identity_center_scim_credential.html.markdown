---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_identity_center_scim_credential"
sidebar_current: "docs-tencentcloud-resource-identity_center_scim_credential"
description: |-
  Provides a resource to create an identity center scim credential
---

# tencentcloud_identity_center_scim_credential

Provides a resource to create an identity center scim credential

## Example Usage

```hcl
resource "tencentcloud_identity_center_scim_credential" "identity_center_scim_credential" {
  zone_id = "z-xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String, ForceNew) Space ID. z-prefix starts with 12 random digits/lowercase letters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - SCIM create time.
* `credential_id` - SCIM key ID. scimcred-prefix and followed by 12 random digits/lowercase letters.
* `credential_type` - SCIM credential type.
* `expire_time` - SCIM expire time.
* `status` - SCIM key status, Enabled-On, Disabled-Closed.


## Import

organization identity_center_scim_credential can be imported using the id, e.g.

```
terraform import tencentcloud_identity_center_scim_credential.identity_center_scim_credential ${zone_id}#${credential_id}
```

