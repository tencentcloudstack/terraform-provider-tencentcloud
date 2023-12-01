---
subcategory: "DNSPOD"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dnspod_custom_line"
sidebar_current: "docs-tencentcloud-resource-dnspod_custom_line"
description: |-
  Provides a resource to create a dnspod custom_line
---

# tencentcloud_dnspod_custom_line

Provides a resource to create a dnspod custom_line

~> **NOTE:** Terraform uses the combined id of doamin and name when importing. When the name changes, the combined id will also change.

## Example Usage

```hcl
resource "tencentcloud_dnspod_custom_line" "custom_line" {
  domain = "dnspod.com"
  name   = "testline8"
  area   = "6.6.6.1-6.6.6.2"
}
```

## Argument Reference

The following arguments are supported:

* `area` - (Required, String) The IP segment of custom line, split with `-`.
* `domain` - (Required, String, ForceNew) Domain.
* `name` - (Required, String) The Name of custom line.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

dnspod custom_line can be imported using the id, e.g.

```
terraform import tencentcloud_dnspod_custom_line.custom_line domain#name
```

