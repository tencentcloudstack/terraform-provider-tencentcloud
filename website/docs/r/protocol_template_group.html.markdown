---
subcategory: "VPC"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_protocol_template_group"
sidebar_current: "docs-tencentcloud-resource-protocol_template_group"
description: |-
  Provides a resource to manage protocol template group.
---

# tencentcloud_protocol_template_group

Provides a resource to manage protocol template group.

## Example Usage

```hcl
resource "tencentcloud_protocol_template_group" "foo" {
  name      = "group-test"
  protocols = ["ipl-axaf24151", "ipl-axaf24152"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) Name of the protocol template group.
* `template_ids` - (Required) Service template ID list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Protocol template group can be imported using the protocol template, e.g.

```
$ terraform import tencentcloud_protocol_template_group.foo ppmg-0np3u974
```

