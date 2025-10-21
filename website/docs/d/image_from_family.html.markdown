---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_image_from_family"
sidebar_current: "docs-tencentcloud-datasource-image_from_family"
description: |-
  Provides query image from family.
---

# tencentcloud_image_from_family

Provides query image from family.

## Example Usage

```hcl
data "tencentcloud_image_from_family" "example" {
  image_family = "business-daily-update"
}
```

## Argument Reference

The following arguments are supported:

* `image_family` - (Required, String) Image family name.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `image` - Information of Image.


