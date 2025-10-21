---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cvm_launch_template_default_version"
sidebar_current: "docs-tencentcloud-resource-cvm_launch_template_default_version"
description: |-
  Provides a resource to create a cvm launch_template_default_version
---

# tencentcloud_cvm_launch_template_default_version

Provides a resource to create a cvm launch_template_default_version

## Example Usage

```hcl
resource "tencentcloud_cvm_launch_template_default_version" "launch_template_default_version" {
  launch_template_id = "lt-34vaef8fe"
  default_version    = 2
}
```

## Argument Reference

The following arguments are supported:

* `default_version` - (Required, Int) The number of the version that you want to set as the default version.
* `launch_template_id` - (Required, String, ForceNew) Instance launch template ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cvm launch_template_default_version can be imported using the id, e.g.

```
terraform import tencentcloud_cvm_launch_template_default_version.launch_template_default_version launch_template_id
```

