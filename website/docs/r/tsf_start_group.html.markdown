---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_start_group"
sidebar_current: "docs-tencentcloud-resource-tsf_start_group"
description: |-
  Provides a resource to create a tsf start_group
---

# tencentcloud_tsf_start_group

Provides a resource to create a tsf start_group

## Example Usage

```hcl
resource "tencentcloud_tsf_start_group" "start_group" {
  group_id = "group-ynd95rea"
  operate  = "start"
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, String) GroupId.
* `operate` - (Required, String) Operation, `start`- start the group, `stop`- stop the group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



