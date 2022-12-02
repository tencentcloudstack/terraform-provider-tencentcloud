---
subcategory: "Cloud Audit(Audit)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_audit_track"
sidebar_current: "docs-tencentcloud-resource-audit_track"
description: |-
  Provides a resource to create a audit track
---

# tencentcloud_audit_track

Provides a resource to create a audit track

## Example Usage

```hcl
resource "tencentcloud_audit_track" "track" {
  action_type = "Read"
  event_names = [
    "*",
  ]
  name                  = "terraform_track"
  resource_type         = "*"
  status                = 1
  track_for_all_members = 0

  storage {
    storage_name   = "db90b92c-91d2-46b0-94ac-debbbb21dc4e"
    storage_prefix = "cloudaudit"
    storage_region = "ap-guangzhou"
    storage_type   = "cls"
  }
}
```

## Argument Reference

The following arguments are supported:

* `action_type` - (Required, String) Track interface type, optional:- `Read`: Read interface- `Write`: Write interface- `*`: All interface.
* `event_names` - (Required, Set: [`String`]) Track interface name list:- when ResourceType is `*`, EventNames is must `[&amp;quot;*&amp;quot;]`- when ResourceType is a single product, EventNames support all interface:`[&amp;quot;*&amp;quot;]`- when ResourceType is a single product, EventNames support some interface, up to 10.
* `name` - (Required, String) Track name.
* `resource_type` - (Required, String) Track product, optional:- `*`: All product- Single product, such as `cos`.
* `status` - (Required, Int) Track status, optional:- `0`: Close- `1`: Open.
* `storage` - (Required, List) Track Storage, support `cos` and `cls`.
* `track_for_all_members` - (Optional, Int) Whether to enable the delivery of group member operation logs to the group management account or trusted service management account, optional:- `0`: Close- `1`: Open.

The `storage` object supports the following:

* `storage_name` - (Required, String) Track Storage name:- when StorageType is `cls`, StorageName is cls topicId- when StorageType is `cos`, StorageName is cos bucket name that does not contain `-APPID`.
* `storage_prefix` - (Required, String) Storage path prefix.
* `storage_region` - (Required, String) Storage region.
* `storage_type` - (Required, String) Track Storage type, optional:- `cos`- `cls`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Track create time.


## Import

audit track can be imported using the id, e.g.
```
$ terraform import tencentcloud_audit_track.track track_id
```

