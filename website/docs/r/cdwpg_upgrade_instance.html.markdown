---
subcategory: "CDWPG"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cdwpg_upgrade_instance"
sidebar_current: "docs-tencentcloud-resource-cdwpg_upgrade_instance"
description: |-
  Provides a resource to create a cdwpg cdwpg_upgrade_instance
---

# tencentcloud_cdwpg_upgrade_instance

Provides a resource to create a cdwpg cdwpg_upgrade_instance

## Example Usage

```hcl
resource "tencentcloud_cdwpg_upgrade_instance" "cdwpg_upgrade_instance" {
  instance_id     = "cdwpg-zpiemnyd"
  package_version = "3.16.9.4"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance id.
* `package_version` - (Required, String, ForceNew) Package version.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



