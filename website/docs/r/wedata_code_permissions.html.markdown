---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_code_permissions"
sidebar_current: "docs-tencentcloud-resource-wedata_code_permissions"
description: |-
  Provides a resource to create a WeData code permissions
---

# tencentcloud_wedata_code_permissions

Provides a resource to create a WeData code permissions

## Example Usage

```hcl
resource "tencentcloud_wedata_code_permissions" "example" {
  project_id = "3108707295180644352"
  authorize_permission_objects {
    resource {
      resource_type        = ""
      resource_id          = ""
      resource_id_for_path = ""
      resource_cfs_path    = ""
    }

    authorize_subjects {
      subject_type   = ""
      subject_values = []
      privileges     = []
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `authorize_permission_objects` - (Required, List) Permission operation objects.
* `project_id` - (Required, String, ForceNew) Project ID.

The `authorize_permission_objects` object supports the following:

* `authorize_subjects` - (Optional, List) Authorization details.
* `resource` - (Optional, List) Authorization resource information, including resourceId and resourceType.

The `authorize_subjects` object of `authorize_permission_objects` supports the following:

* `privileges` - (Optional, Set) Permission list.
* `subject_type` - (Optional, String) Subject type (user: user, role: role, group: group).
* `subject_values` - (Optional, Set) Subject value list.

The `resource` object of `authorize_permission_objects` supports the following:

* `resource_cfs_path` - (Optional, String) CFS path.
* `resource_id_for_path` - (Optional, String) Full ID path, used for recursive authentication.
* `resource_id` - (Optional, String) Resource ID: directory ID or script ID.
* `resource_type` - (Optional, String) Resource type, can only be these two types: folder, script.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateCodePermissions](https://www.tencentcloud.com/document/api/1174/77907)
* Read: [ListCodePermissions](https://www.tencentcloud.com/document/api/1174/77902)
* Delete: [DeleteCodePermissions](https://www.tencentcloud.com/document/api/1174/77906)

## Import

WeData code permissions can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_code_permissions.example wedata_code_permissions_id
```

