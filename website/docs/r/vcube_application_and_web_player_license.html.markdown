---
subcategory: "VCube"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vcube_application_and_web_player_license"
sidebar_current: "docs-tencentcloud-resource-vcube_application_and_web_player_license"
description: |-
  Provides a resource to create a VCube application and web player license
---

# tencentcloud_vcube_application_and_web_player_license

Provides a resource to create a VCube application and web player license

## Example Usage

```hcl
resource "tencentcloud_vcube_application_and_web_player_license" "example" {
  app_name = "tf-example"
  domain_list = [
    "www.example1.com",
    "www.example2.com",
    "www.example3.com",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `app_name` - (Required, String, ForceNew) Application name.
* `domain_list` - (Required, Set: [`String`], ForceNew) Domain list.

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

VCube application and web player license can be imported using the id, e.g.

```
terraform import tencentcloud_vcube_application_and_web_player_license.example 1513
```

