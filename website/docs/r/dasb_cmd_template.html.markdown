---
subcategory: "Bastion Host(BH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dasb_cmd_template"
sidebar_current: "docs-tencentcloud-resource-dasb_cmd_template"
description: |-
  Provides a resource to create a dasb cmd_template
---

# tencentcloud_dasb_cmd_template

Provides a resource to create a dasb cmd_template

## Example Usage

```hcl
resource "tencentcloud_dasb_cmd_template" "example" {
  name     = "tf_example"
  cmd_list = "rm -rf*"
}
```

## Argument Reference

The following arguments are supported:

* `cmd_list` - (Required, String) Command list, n separated, maximum length 32768 bytes.
* `name` - (Required, String) Template name, maximum length 32 characters, cannot contain blank characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

dasb cmd_template can be imported using the id, e.g.

```
terraform import tencentcloud_dasb_cmd_template.example 15
```

