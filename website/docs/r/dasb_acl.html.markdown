---
subcategory: "Bastion Host(BH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dasb_acl"
sidebar_current: "docs-tencentcloud-resource-dasb_acl"
description: |-
  Provides a resource to create a dasb acl
---

# tencentcloud_dasb_acl

Provides a resource to create a dasb acl

## Example Usage

```hcl
resource "tencentcloud_dasb_acl" "example" {
  name                    = "tf_example"
  allow_disk_redirect     = true
  allow_any_account       = false
  allow_clip_file_up      = true
  allow_clip_file_down    = true
  allow_clip_text_up      = true
  allow_clip_text_down    = true
  allow_file_up           = true
  allow_file_down         = true
  max_file_up_size        = 0
  max_file_down_size      = 0
  user_id_set             = ["6", "2"]
  user_group_id_set       = ["6", "36"]
  device_id_set           = ["39", "81"]
  device_group_id_set     = ["2", "3"]
  account_set             = ["root"]
  cmd_template_id_set     = ["1", "7"]
  ac_template_id_set      = []
  allow_disk_file_up      = true
  allow_disk_file_down    = true
  allow_shell_file_up     = true
  allow_shell_file_down   = true
  allow_file_del          = true
  allow_access_credential = true
  department_id           = "1.2"
  validate_from           = "2023-09-22T00:00:00+08:00"
  validate_to             = "2024-09-23T00:00:00+08:00"
}
```

## Argument Reference

The following arguments are supported:

* `allow_any_account` - (Required, Bool) Allow any account.
* `allow_disk_redirect` - (Required, Bool) Allow disk redirect.
* `name` - (Required, String) Acl name.
* `ac_template_id_set` - (Optional, Set: [`String`]) Associate high-risk DB template IDs.
* `account_set` - (Optional, Set: [`String`]) Associated accounts.
* `allow_access_credential` - (Optional, Bool) Allow access credential,default allow.
* `allow_clip_file_down` - (Optional, Bool) Allow clip file down.
* `allow_clip_file_up` - (Optional, Bool) Allow clip file up.
* `allow_clip_text_down` - (Optional, Bool) Allow clip text down.
* `allow_clip_text_up` - (Optional, Bool) Allow clip text up.
* `allow_disk_file_down` - (Optional, Bool) Allow disk file download.
* `allow_disk_file_up` - (Optional, Bool) Allow disk file upload.
* `allow_file_del` - (Optional, Bool) Allow sftp file delete.
* `allow_file_down` - (Optional, Bool) Allow sftp file download.
* `allow_file_up` - (Optional, Bool) Allow sftp up file.
* `allow_shell_file_down` - (Optional, Bool) Allow shell file download.
* `allow_shell_file_up` - (Optional, Bool) Allow shell file upload.
* `cmd_template_id_set` - (Optional, Set: [`Int`]) Associated high-risk command template ID.
* `department_id` - (Optional, String) Department id.
* `device_group_id_set` - (Optional, Set: [`Int`]) Associated device group ID.
* `device_id_set` - (Optional, Set: [`Int`]) Associated collection of device IDs.
* `max_file_down_size` - (Optional, Int) File transfer download size limit (reserved parameter, currently unused).
* `max_file_up_size` - (Optional, Int) File upload transfer size limit (artifact parameter, currently unused).
* `user_group_id_set` - (Optional, Set: [`Int`]) Associated user group ID.
* `user_id_set` - (Optional, Set: [`Int`]) Associated set of user IDs.
* `validate_from` - (Optional, String) Access permission effective time, such as: 2021-09-22T00:00:00+08:00If the effective and expiry time are not filled in, the access rights will be valid for a long time.
* `validate_to` - (Optional, String) Access permission expiration time, such as: 2021-09-23T00:00:00+08:00If the effective and expiry time are not filled in, the access rights will be valid for a long time.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

dasb acl can be imported using the id, e.g.

```
terraform import tencentcloud_dasb_acl.example 132
```

