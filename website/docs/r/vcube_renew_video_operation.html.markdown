---
subcategory: "VCube"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vcube_renew_video_operation"
sidebar_current: "docs-tencentcloud-resource-vcube_renew_video_operation"
description: |-
  Provides a resource to create a VCube renew video operation
---

# tencentcloud_vcube_renew_video_operation

Provides a resource to create a VCube renew video operation

~> **NOTE:** Resource `tencentcloud_vcube_renew_video_operation` can be directly invoked to renew the license within 30 days before its expiration. Once the renewal is successful, an additional year will be added immediately.

## Example Usage

```hcl
resource "tencentcloud_vcube_renew_video_operation" "example" {
  license_id = 1513
}
```

## Argument Reference

The following arguments are supported:

* `license_id` - (Required, Int, ForceNew) License ID for video playback renewal.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



