---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_service_linked_role"
sidebar_current: "docs-tencentcloud-resource-cam_service_linked_role"
description: |-
  Provides a resource to create a cam service_linked_role
---

# tencentcloud_cam_service_linked_role

Provides a resource to create a cam service_linked_role

## Example Usage

```hcl
resource "tencentcloud_cam_service_linked_role" "service_linked_role" {
  qcs_service_name = ["cvm.qcloud.com", "ekslog.tke.cloud.tencent.com"]
  custom_suffix    = "tf"
  description      = "desc cam"
  tags = {
    "createdBy" = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `qcs_service_name` - (Required, Set: [`String`], ForceNew) Authorization service, the Tencent Cloud service principal with this role attached.
* `custom_suffix` - (Optional, String, ForceNew) The custom suffix, based on the string you provide, is combined with the prefix provided by the service to form the full role name. This field is not allowed to contain the character `_`.
* `description` - (Optional, String) role description.
* `tags` - (Optional, Map) Tag description list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



