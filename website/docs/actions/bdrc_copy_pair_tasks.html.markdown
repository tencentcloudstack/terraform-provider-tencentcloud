---
subcategory: "Provider Meta"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_bdrc_copy_pair_tasks"
sidebar_current: "docs-tencentcloud-action-bdrc_copy_pair_tasks"
description: |-
  Provides an action to launch a group of BDRC (Business Data Resilient Disaster Recovery) copy pair tasks via the RunCopyPairTasks API. This is a one-time operation; no cloud-side state is persisted after the action completes.
---

# tencentcloud_bdrc_copy_pair_tasks

Provides an action to launch a group of BDRC (Business Data Resilient Disaster Recovery) copy pair tasks via the RunCopyPairTasks API. This is a one-time operation; no cloud-side state is persisted after the action completes.

~> **NOTE:** Actions are supported in HashiCorp Terraform version 1.14 and later.

## Example Usage

```hcl
action "tencentcloud_bdrc_copy_pair_tasks" "example" {
  config {
    copy_pair_ids = [
      "pair-1111222233334444",
      "pair-5555666677778888",
    ]
    copy_pair_type = "DISK"
  }
}
```

## Argument Reference

The following arguments are supported:

* `copy_pair_ids` - (Required, List) List of copy pair IDs to launch.
* `copy_pair_type` - (Required, String) Type of the copy pairs to launch. Valid values: DISK, INSTANCE, CFS.


