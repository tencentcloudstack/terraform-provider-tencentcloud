---
subcategory: "Performance Testing Service(PTS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_pts_project"
sidebar_current: "docs-tencentcloud-resource-pts_project"
description: |-
  Provides a resource to create a pts project
---

# tencentcloud_pts_project

Provides a resource to create a pts project

## Example Usage

```hcl
resource "tencentcloud_pts_project" "project" {
  name        = "ptsObjectName-1"
  description = "desc"
  tags {
    tag_key   = "createdBy"
    tag_value = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) ptsObjectName, which must be required.
* `description` - (Optional, String) Pts object description.
* `tags` - (Optional, List) Tags List.

The `tags` object supports the following:

* `tag_key` - (Optional, String) tag key.
* `tag_value` - (Optional, String) tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `app_id` - App ID.
* `created_at` - Creation time.
* `status` - Project status.
* `sub_account_uin` - Sub-user ID.
* `uin` - User ID.
* `updated_at` - Update time.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateProject](https://cloud.tencent.com/document/api/1484/78092) - This API is used to create a project.
* Read: [DescribeProjects](https://cloud.tencent.com/document/api/1484/78074) - The API is used to query project list.
* Update: [UpdateProject](https://cloud.tencent.com/document/api/1484/78058) - This API is used to update a project.
* Delete: [DeleteProjects](https://cloud.tencent.com/document/api/1484/78086) - This API is used to delete projects.

## Import

pts project can be imported using the id, e.g.
```
$ terraform import tencentcloud_pts_project.project project-1ep27k1m
```

