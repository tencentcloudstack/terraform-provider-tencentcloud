---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_quit_organization_operation"
sidebar_current: "docs-tencentcloud-resource-organization_quit_organization_operation"
description: |-
  Provides a resource to create a organization quit_organization_operation
---

# tencentcloud_organization_quit_organization_operation

Provides a resource to create a organization quit_organization_operation

## Example Usage

```hcl
resource "tencentcloud_organization_quit_organization_operation" "quit_organization_operation" {
  org_id = 45155
}
```

## Argument Reference

The following arguments are supported:

* `org_id` - (Required, Int, ForceNew) Organization ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

organization quit_organization_operation can be imported using the id, e.g.

```
terraform import tencentcloud_organization_quit_organization_operation.quit_organization_operation quit_organization_operation_id
```

