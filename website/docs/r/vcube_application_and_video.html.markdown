---
subcategory: "VCube"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vcube_application_and_video"
sidebar_current: "docs-tencentcloud-resource-vcube_application_and_video"
description: |-
  Provides a resource to create a VCube application and video
---

# tencentcloud_vcube_application_and_video

Provides a resource to create a VCube application and video

## Example Usage

```hcl
resource "tencentcloud_vcube_application_and_video" "example" {
  app_name  = "tf-example"
  bundle_id = "com.example.appName"
}
```

### Or

```hcl
resource "tencentcloud_vcube_application_and_video" "example" {
  app_name     = "tf-example"
  package_name = "com.example.appName"
}
```

### Or

```hcl
resource "tencentcloud_vcube_application_and_video" "example" {
  app_name     = "tf-example"
  bundle_id    = "com.example.appName"
  package_name = "com.example.appName"
}
```

## Argument Reference

The following arguments are supported:

* `app_name` - (Required, String, ForceNew) Application name.
* `bundle_id` - (Optional, String, ForceNew) IOS bundle ID. Choose at least one of `bundle_id` and `package_name`.
* `package_name` - (Optional, String, ForceNew) Android package name. Choose at least one of `bundle_id` and `package_name`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `app_id` - Account App ID.
* `app_type` - Application type: formal: formal application, test: test application.
* `application_id` - User Application ID.
* `license_id` - License ID.
* `license_key` - License key.
* `license_url` - License url.


## Import

VCube application and video can be imported using the id, e.g.

```
terraform import tencentcloud_vcube_application_and_video.example 1509
```

