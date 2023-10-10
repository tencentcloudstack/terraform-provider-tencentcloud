---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_access_key"
sidebar_current: "docs-tencentcloud-resource-cam_access_key"
description: |-
  Provides a resource to create a cam access_key
---

# tencentcloud_cam_access_key

Provides a resource to create a cam access_key

## Example Usage

```hcl
resource "tencentcloud_cam_access_key" "access_key" {
  target_uin = 100033690181
}
```

### Update

```hcl
resource "tencentcloud_cam_access_key" "access_key" {
  target_uin = 100033690181
  status     = "Inactive"
}
```

### Encrypted

```hcl
resource "tencentcloud_cam_access_key" "access_key" {
  target_uin = 100033690181
  pgp_key    = "keybase:some_person_that_exists"
}
```

## Argument Reference

The following arguments are supported:

* `access_key` - (Optional, String) Access_key is the access key identification, required when updating.
* `pgp_key` - (Optional, String, ForceNew) Either a base-64 encoded PGP public key, or a keybase username in the form keybase:some_person_that_exists, for use in the encrypted_secret output attribute. If providing a base-64 encoded PGP public key, make sure to provide the "raw" version and not the "armored" one (e.g. avoid passing the -a option to gpg --export).
* `status` - (Optional, String) Key status, activated (Active) or inactive (Inactive), required when updating.
* `target_uin` - (Optional, Int) Specify user Uin, if not filled, the access key is created for the current user by default.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `encrypted_secret_access_key` - Encrypted secret, base64 encoded, if pgp_key was specified. This attribute is not available for imported resources. The encrypted secret may be decrypted using the command line, for example: terraform output -raw encrypted_secret | base64 --decode | keybase pgp decrypt.
* `key_fingerprint` - Fingerprint of the PGP key used to encrypt the secret. This attribute is not available for imported resources.
* `secret_access_key` - Access key (key is only visible when created, please keep it properly).


## Import

cam access_key can be imported using the id, e.g.

```
terraform import tencentcloud_cam_access_key.access_key access_key_id
```

