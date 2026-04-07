---
subcategory: "Bastion Host(BH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_bh_acl"
sidebar_current: "docs-tencentcloud-resource-bh_acl"
description: |-
  Provides a BH (Bastion Host) acl.
---

# tencentcloud_bh_acl

Provides a BH (Bastion Host) acl.

## Example Usage

```hcl
resource "tencentcloud_bh_acl" "example" {
  name                    = "tf-example"
  allow_disk_redirect     = true
  allow_any_account       = false
  allow_clip_file_up      = true
  allow_clip_file_down    = false
  allow_clip_text_up      = false
  allow_clip_text_down    = true
  allow_file_up           = false
  allow_file_down         = true
  allow_disk_file_up      = true
  allow_disk_file_down    = false
  allow_shell_file_up     = false
  allow_shell_file_down   = false
  allow_file_del          = false
  allow_access_credential = true
  allow_keyboard_logger   = false
}
```

## Argument Reference

The following arguments are supported:

* `allow_any_account` - (Required, Bool) Whether to allow any account to log in.
* `allow_disk_redirect` - (Required, Bool) Whether to enable disk mapping.
* `name` - (Required, String) Access permission name, maximum 32 characters, cannot contain whitespace characters.
* `ac_template_id_set` - (Optional, Set: [`String`]) Associated high-risk DB template ID set.
* `account_set` - (Optional, Set: [`String`]) Associated account set.
* `allow_access_credential` - (Optional, Bool) Whether to allow the use of access credentials. Default is allowed.
* `allow_clip_file_down` - (Optional, Bool) Whether to enable clipboard file download.
* `allow_clip_file_up` - (Optional, Bool) Whether to enable clipboard file upload.
* `allow_clip_text_down` - (Optional, Bool) Whether to enable clipboard text (including images) download.
* `allow_clip_text_up` - (Optional, Bool) Whether to enable clipboard text (including images) upload.
* `allow_disk_file_down` - (Optional, Bool) Whether to enable RDP disk mapping file download.
* `allow_disk_file_up` - (Optional, Bool) Whether to enable RDP disk mapping file upload.
* `allow_file_del` - (Optional, Bool) Whether to enable SFTP file deletion.
* `allow_file_down` - (Optional, Bool) Whether to enable SFTP file download.
* `allow_file_up` - (Optional, Bool) Whether to enable SFTP file upload.
* `allow_keyboard_logger` - (Optional, Bool) Whether to allow keyboard logging.
* `allow_shell_file_down` - (Optional, Bool) Whether to enable rz sz file download.
* `allow_shell_file_up` - (Optional, Bool) Whether to enable rz sz file upload.
* `app_asset_id_set` - (Optional, Set: [`Int`]) Associated application asset ID set.
* `cmd_template_id_set` - (Optional, Set: [`Int`]) Associated high-risk command template ID set.
* `department_id` - (Optional, String) Department ID to which the access permission belongs, e.g.: `1.2.3`.
* `device_group_id_set` - (Optional, Set: [`Int`]) Associated asset group ID set.
* `device_id_set` - (Optional, Set: [`Int`]) Associated asset ID set.
* `max_access_credential_duration` - (Optional, Int) Maximum validity period of access credentials (in seconds). Must be a multiple of 86400 when access credentials are enabled.
* `max_file_down_size` - (Optional, Int) File transfer download size limit (reserved parameter, not currently used).
* `max_file_up_size` - (Optional, Int) File transfer upload size limit (reserved parameter, not currently used).
* `user_group_id_set` - (Optional, Set: [`Int`]) Associated user group ID set.
* `user_id_set` - (Optional, Set: [`Int`]) Associated user ID set.
* `validate_from` - (Optional, String) Access permission effective time in ISO8601 format, e.g.: `2021-09-22T00:00:00+00:00`. If not set, the permission is permanently valid.
* `validate_to` - (Optional, String) Access permission expiration time in ISO8601 format, e.g.: `2021-09-23T00:00:00+00:00`. If not set, the permission is permanently valid.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `acl_id` - Access permission ID.


## Import

BH acl can be imported using the id, e.g.

```
terraform import tencentcloud_bh_acl.example 1374
```

