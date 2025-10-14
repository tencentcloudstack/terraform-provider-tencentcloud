---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_ops_task_owner"
sidebar_current: "docs-tencentcloud-resource-wedata_ops_task_owner"
description: |-
  Provides a resource to create a wedata ops task owner
---

# tencentcloud_wedata_ops_task_owner

Provides a resource to create a wedata ops task owner

## Example Usage

```hcl
resource "tencentcloud_wedata_ops_task_owner" "wedata_ops_task_owner" {
  owner_uin  = "100029411056;100042282926"
  project_id = "2430455587205529600"
  task_id    = "20251009144419600"
}
```

## Argument Reference

The following arguments are supported:

* `owner_uin` - (Required, String) Task Owner ID. For multiple owners, separate them with `;`, for example: `100029411056;100042282926`.
* `project_id` - (Required, String) Project id.
* `task_id` - (Required, String) Task id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

wedata ops task owner can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_ops_task_owner.wedata_ops_task_owner projectId#askId
```

