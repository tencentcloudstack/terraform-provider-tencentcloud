---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_org_share_unit"
sidebar_current: "docs-tencentcloud-resource-organization_org_share_unit"
description: |-
  Provides a resource to create a organization org_share_unit
---

# tencentcloud_organization_org_share_unit

Provides a resource to create a organization org_share_unit

## Example Usage

```hcl
resource "tencentcloud_organization_org_share_unit" "org_share_unit" {
  name        = "iac-test"
  area        = "ap-guangzhou"
  description = "iac-test"
}
```

## Argument Reference

The following arguments are supported:

* `area` - (Required, String) Shared unit region. The regions that support sharing can be obtained through the DescribeShareAreas interface.
* `name` - (Required, String) Shared unit name. It only supports a combination of uppercase and lowercase letters, numbers, -, and _, with a length of 3-128 characters.
* `description` - (Optional, String) Shared unit description. Up to 128 characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `unit_id` - Shared unit region. The regions that support sharing can be obtained through the DescribeShareAreas interface.


## Import

organization org_share_unit can be imported using the id, e.g.

```
terraform import tencentcloud_organization_org_share_unit.org_share_unit area#unit_id
```

