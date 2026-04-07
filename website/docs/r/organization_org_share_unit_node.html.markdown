---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_org_share_unit_node"
sidebar_current: "docs-tencentcloud-resource-organization_org_share_unit_node"
description: |-
  Provides a resource to create an organization org share unit node
---

# tencentcloud_organization_org_share_unit_node

Provides a resource to create an organization org share unit node

## Example Usage

```hcl
resource "tencentcloud_organization_org_share_unit_node" "example" {
  unit_id = "us-xxxxx"
  node_id = 123456
}
```

## Argument Reference

The following arguments are supported:

* `node_id` - (Required, Int, ForceNew) Organization department ID.
* `unit_id` - (Required, String, ForceNew) Shared unit ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

organization org share unit node can be imported using the id, e.g.

```
terraform import tencentcloud_organization_org_share_unit_node.example us-xxxxx#123456
```

