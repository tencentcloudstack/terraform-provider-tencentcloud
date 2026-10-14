---
subcategory: "Video on Demand(VOD)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vod_sub_application"
sidebar_current: "docs-tencentcloud-resource-vod_sub_application"
description: |-
  Provide a resource to create a VOD sub application.
---

# tencentcloud_vod_sub_application

Provide a resource to create a VOD sub application.

## Example Usage

### ### Basic Usage

```hcl
resource "tencentcloud_vod_sub_application" "example" {
  name        = "tf-example"
  status      = "On"
  description = "this is sub application"
}
```

### ### Tags Update Example

```hcl
resource "tencentcloud_vod_sub_application" "example" {
  name        = "tf-example"
  status      = "On"
  description = "Sub application with updatable tags"

  tags = {
    "team"        = "media"
    "environment" = "production"
  }
}

# Tags can be updated without recreating the resource
# Modify the tags map to add, update, or remove tags
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Sub application name, which can contain up to 64 letters, digits, underscores, and hyphens (such as test_ABC-123) and must be unique under a user.
* `status` - (Required, String) Sub appliaction status.
* `description` - (Optional, String) Sub application description.
* `tags` - (Optional, Map) Tag key-value pairs for resource management. Maximum 10 tags.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - The time when the sub application was created.
* `sub_app_id` - Sub application ID.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateSubAppId](https://www.tencentcloud.com/document/product/266/37569) - This API is used to create a VOD application.
* Create: [DescribeSubAppIds](https://www.tencentcloud.com/document/product/266/34177) - This API is used to get the application list of the current account.
* Create: [ModifySubAppIdStatus](https://www.tencentcloud.com/document/product/266/34173) - This API is used to enable or disable an application. A disabled application will have its corresponding domain blocked and console access restricted.
* Read: [DescribeSubAppIds](https://www.tencentcloud.com/document/product/266/34177) - This API is used to get the application list of the current account.
* Update: [DescribeSubAppIds](https://www.tencentcloud.com/document/product/266/34177) - This API is used to get the application list of the current account.
* Update: [ModifyResourceTags](https://www.tencentcloud.com/document/product/651/33012) - This API is used to modify all tags associated with a resource.
* Update: [ModifySubAppIdInfo](https://www.tencentcloud.com/document/product/266/34174) - This API is used to change application information, but the default application info is not allowed to be modified.
* Update: [ModifySubAppIdStatus](https://www.tencentcloud.com/document/product/266/34173) - This API is used to enable or disable an application. A disabled application will have its corresponding domain blocked and console access restricted.
* Delete: [ModifySubAppIdStatus](https://www.tencentcloud.com/document/product/266/34173) - This API is used to enable or disable an application. A disabled application will have its corresponding domain blocked and console access restricted.

## Import

VOD sub application can be imported using the name and id separated by `name#sub_app_id`, e.g.

```
terraform import tencentcloud_vod_sub_application.example tf-example#1500066377
```

