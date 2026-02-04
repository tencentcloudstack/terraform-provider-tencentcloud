---
subcategory: "Cloud Audit(Audit)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_events_audit_track"
sidebar_current: "docs-tencentcloud-resource-events_audit_track"
description: |-
  Provides a resource to create events audit track
---

# tencentcloud_events_audit_track

Provides a resource to create events audit track

## Example Usage

```hcl
resource "tencentcloud_events_audit_track" "example" {
  name = "track_example"

  status                = 1
  track_for_all_members = 0

  storage {
    storage_name   = "393953ac-5c1b-457d-911d-376271b1b4f2"
    storage_prefix = "cloudaudit"
    storage_region = "ap-guangzhou"
    storage_type   = "cls"
  }

  filters {
    resource_fields {
      resource_type = "cam"
      action_type   = "*"
      event_names   = ["AddSubAccount", "AddSubAccountCheckingMFA"]
    }
    resource_fields {
      resource_type = "cvm"
      action_type   = "*"
      event_names   = ["*"]
    }
    resource_fields {
      resource_type = "tke"
      action_type   = "*"
      event_names   = ["*"]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Required, List) Data filtering criteria.
* `name` - (Required, String, ForceNew) Tracking set name, which can only contain 3-48 letters, digits, hyphens, and underscores.
* `status` - (Required, Int) Tracking set status (0: Not enabled; 1: Enabled).
* `storage` - (Required, List) Storage type of shipped data. Valid values: `cos`, `cls` and `ckafka`.
* `track_for_all_members` - (Optional, Int) Whether to enable the feature of shipping organization members operation logs to the organization admin account or the trusted service admin account (0: Not enabled; 1: Enabled. This feature can only be enabled by the organization admin account or the trusted service admin account).

The `filters` object supports the following:

* `resource_fields` - (Optional, List) Resource filtering conditions.

The `resource_fields` object of `filters` supports the following:

* `action_type` - (Required, String) Tracking set event type (`Read`: Read; `Write`: Write; `*`: All).
* `event_names` - (Required, Set) The list of API names of tracking set events. When `ResourceType` is `*`, the value of `EventNames` must be `*`. When `ResourceType` is a specified product, the value of `EventNames` can be `*`. When `ResourceType` is `cos` or `cls`, up to 10 APIs are supported.
* `resource_type` - (Required, String) The product to which the tracking set event belongs. The value can be a single product such as `cos`, or `*` that indicates all products.

The `storage` object supports the following:

* `storage_name` - (Required, String) Storage name. For COS, the storage name is the custom bucket name, which can contain up to 50 lowercase letters, digits, and hyphens. It cannot contain "-APPID" and cannot start or end with a hyphen. For CLS, the storage name is the log topic ID, which can contain 1-50 characters.
* `storage_prefix` - (Required, String) Storage directory prefix. The COS log file prefix can only contain 3-40 letters and digits.
* `storage_region` - (Required, String) StorageRegion *string `json:'StorageRegion,omitnil,omitempty' name: 'StorageRegion'`.
* `storage_type` - (Required, String) Storage type (Valid values: cos, cls, ckafka).
* `storage_account_id` - (Optional, String) Designated to store user ID.
* `storage_app_id` - (Optional, String) Designated to store user app ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `track_id` - Whether the log list has come to an end. `true`: Yes. Pagination is not required.


## Import

events audit track can be imported using the id, e.g.
```
$ terraform import tencentcloud_events_audit_track.example 24283
```

