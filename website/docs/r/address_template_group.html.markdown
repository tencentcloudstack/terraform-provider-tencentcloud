---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_address_template_group"
sidebar_current: "docs-tencentcloud-resource-address_template_group"
description: |-
  Provides a resource to manage address template group.
---

# tencentcloud_address_template_group

Provides a resource to manage address template group.

## Example Usage

```hcl
resource "tencentcloud_address_template_group" "foo" {
  name      = "group-test"
  addresses = ["ipl-axaf24151", "ipl-axaf24152"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) Name of the address template group.
* `template_ids` - (Required) Template ID list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Address template group can be imported using the id, e.g.

```
$ terraform import tencentcloud_address_template_group.foo ipmg-0np3u974
```

