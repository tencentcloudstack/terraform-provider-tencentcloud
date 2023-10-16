---
subcategory: "Key Management Service(KMS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kms_overwrite_white_box_device_fingerprints"
sidebar_current: "docs-tencentcloud-resource-kms_overwrite_white_box_device_fingerprints"
description: |-
  Provides a resource to create a kms overwrite_white_box_device_fingerprints
---

# tencentcloud_kms_overwrite_white_box_device_fingerprints

Provides a resource to create a kms overwrite_white_box_device_fingerprints

## Example Usage

```hcl
resource "tencentcloud_kms_overwrite_white_box_device_fingerprints" "example" {
  key_id = "23e80852-1e38-11e9-b129-5cb9019b4b01"
}
```

## Argument Reference

The following arguments are supported:

* `key_id` - (Required, String, ForceNew) CMK unique identifier.
* `device_fingerprints` - (Optional, List, ForceNew) Device fingerprint list.

The `device_fingerprints` object supports the following:

* `identity` - (Required, String) identity.
* `description` - (Optional, String) Description.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



