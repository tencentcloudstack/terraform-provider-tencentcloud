---
subcategory: "Cloud Infinite(CI)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ci_original_image_protection"
sidebar_current: "docs-tencentcloud-resource-ci_original_image_protection"
description: |-
  Manage original image protection functionality
---

# tencentcloud_ci_original_image_protection

Manage original image protection functionality

## Example Usage

```hcl
resource "tencentcloud_ci_original_image_protection" "foo" {
  bucket = "examplebucket-1250000000"
  status = "on"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String, ForceNew) The name of a bucket, the format should be [custom name]-[appid], for example `mycos-1258798060`.
* `status` - (Required, String) Whether original image protection is set, options: on/off.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Resource original image protection can be imported using the id, e.g.

```
$ terraform import tencentcloud_ci_original_image_protection.example examplebucket-1250000000
```

