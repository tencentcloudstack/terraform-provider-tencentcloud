---
subcategory: "Config"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_config_recorder_config"
sidebar_current: "docs-tencentcloud-resource-config_recorder_config"
description: |-
  Provides a resource to manage Config recorder configuration (global singleton).
---

# tencentcloud_config_recorder_config

Provides a resource to manage Config recorder configuration (global singleton).

## Example Usage

### Enable monitoring and specify resource types

```hcl
resource "tencentcloud_config_recorder_config" "example" {
  status = true
  resource_types = [
    "QCS::CAM::Group",
    "QCS::CAM::Role",
    "QCS::CAM::Policy",
    "QCS::CAM::User",
    "QCS::CVM::Instance",
    "QCS::COS::Bucket",
  ]
}
```

### Disable monitoring

```hcl
resource "tencentcloud_config_recorder_config" "example" {
  status = false
}
```

## Argument Reference

The following arguments are supported:

* `status` - (Required, Bool) Whether to enable resource monitoring. true: enable (OpenConfigRecorder), false: disable (CloseConfigRecorder).
* `resource_types` - (Optional, List: [`String`]) Resource type list to monitor (e.g. QCS::CAM::Group, QCS::CVM::Instance).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Recorder creation time.
* `open_count` - Number of times monitoring was opened today.
* `trigger_count` - Number of snapshots taken today.
* `update_count` - Number of monitoring range updates today.


## Import

Config recorder config can be imported using its token ID, e.g.

```
terraform import tencentcloud_config_recorder_config.example <id>
```

