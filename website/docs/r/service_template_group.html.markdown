---
subcategory: "VPC"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_service_template_group"
sidebar_current: "docs-tencentcloud-resource-service_template_group"
description: |-
  Provides a resource to manage service template group.
---

# tencentcloud_service_template_group

Provides a resource to manage service template group.

## Example Usage

```hcl
resource "tencentcloud_service_template_group" "foo" {
  name     = "group-test"
  services = ["ipl-axaf24151", "ipl-axaf24152"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) Name of the service template group.
* `template_ids` - (Required) Service template ID list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

CAM user can be imported using the service template, e.g.

```
$ terraform import tencentcloud_service_template.foo ppmg-0np3u974
```

