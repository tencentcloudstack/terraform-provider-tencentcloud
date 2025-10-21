---
subcategory: "MapReduce(EMR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_emr_yarn"
sidebar_current: "docs-tencentcloud-resource-emr_yarn"
description: |-
  Provides a resource to create a emr emr_yarn
---

# tencentcloud_emr_yarn

Provides a resource to create a emr emr_yarn

## Example Usage

```hcl
resource "tencentcloud_emr_yarn" "emr_yarn" {
  instance_id              = "emr-rzrochgp"
  enable_resource_schedule = true
  scheduler                = "fair"
  fair_global_config {
    user_max_apps_default = 1000
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) EMR Instance ID.
* `capacity_global_config` - (Optional, List) Information about capacity scheduler.
* `disable_resource_schedule_sync` - (Optional, Bool) Whether to synchronize when closing.
* `enable_resource_schedule` - (Optional, Bool) Whether the resource scheduling function is enabled.
* `fair_global_config` - (Optional, List) Information about fair scheduler.
* `scheduler` - (Optional, String) The latest resource scheduler.

The `capacity_global_config` object supports the following:

* `default_settings` - (Optional, Set) Advanced settings.
* `enable_label` - (Optional, Bool) Whether label scheduling is turned on.
* `label_dir` - (Optional, String) If label scheduling is enabled, the path where label information is stored.
* `queue_mapping_override` - (Optional, Bool) Whether to overwrite the user-specified queue. True means overwriting.

The `default_settings` object of `capacity_global_config` supports the following:

* `name` - (Required, String) Name, as the key for the input parameter.
* `value` - (Required, String) Value corresponding to tame.

The `fair_global_config` object supports the following:

* `user_max_apps_default` - (Optional, Int) Corresponding to the page procedural upper limit.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `active_scheduler` - Resource dispatchers currently in effect.


## Import

emr emr_yarn can be imported using the id, e.g.

```
terraform import tencentcloud_emr_yarn.emr_yarn emr_instance_id
```

