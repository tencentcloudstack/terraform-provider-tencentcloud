---
subcategory: "TCR"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_namespace"
sidebar_current: "docs-tencentcloud-resource-tcr_namespace"
description: |-
  Use this resource to create tcr namespace.
---

# tencentcloud_tcr_namespace

Use this resource to create tcr namespace.

## Example Usage

```hcl
resource "tencentcloud_tcr_namespace" "foo" {
  instance_id = ""
  name        = "example"
  is_public   = true
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) Id of the TCR instance.
* `name` - (Required, ForceNew) Name of the TCR namespace. Valid length is 2~30. It can only contain lowercase letters, numbers and separators (`.`, `_`, `-`), and cannot start, end or continue with separators.
* `is_public` - (Optional) Indicate that the namespace is public or not. Default is `false`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tcr namespace can be imported using the id, e.g.

```
$ terraform import tencentcloud_tcr_namespace.foo cls-cda1iex1#namespace
```

