---
subcategory: "Intelligent Global Traffic Manager(IGTM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_igtm_package_task"
sidebar_current: "docs-tencentcloud-resource-igtm_package_task"
description: |-
  Provides a resource to create a IGTM package task
---

# tencentcloud_igtm_package_task

Provides a resource to create a IGTM package task

## Example Usage

```hcl
resource "tencentcloud_igtm_package_task" "example" {
  task_detection_quantity = 100
  auto_renew              = 2
  time_span               = 1
  auto_voucher            = 1
}
```

## Argument Reference

The following arguments are supported:

* `auto_renew` - (Required, Int) Auto renewal: 1 enable auto renewal; 2 disable auto renewal.
* `task_detection_quantity` - (Required, Int) Value range: 1~10000.
* `auto_voucher` - (Optional, Int) Whether to automatically select vouchers, 1 yes; 0 no, default is 0.
* `time_span` - (Optional, Int) Package duration in months, required for creation and renewal. Value range: 1~120.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `task_id` - Task ID.


## Import

IGTM package task can be imported using the id, e.g.

```
terraform import tencentcloud_igtm_package_task.example task-dahygvmzawgn
```

