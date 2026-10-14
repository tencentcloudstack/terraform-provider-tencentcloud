---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_project_member"
sidebar_current: "docs-tencentcloud-resource-wedata_project_member"
description: |-
  Provides a resource to create a WeData project member
---

# tencentcloud_wedata_project_member

Provides a resource to create a WeData project member

~> **NOTE:** This resource must exclusive in one share unit, do not declare additional roleIds resources of this project member elsewhere.

## Example Usage

```hcl
resource "tencentcloud_wedata_project" "example" {
  project {
    project_name  = "tf_example"
    display_name  = "display_name"
    project_model = "SIMPLE"
  }

  dlc_info {
    compute_resources     = ["svmgao_stability"]
    region                = "ap-guangzhou"
    default_database      = "db_name"
    standard_mode_env_tag = "Dev"
    access_account        = "OWNER"
  }

  resource_ids = [
    "20250909193110713075",
    "20250820215449817917"
  ]

  status = 1
}

resource "tencentcloud_wedata_project_member" "example" {
  project_id = tencentcloud_wedata_project.example.id
  user_uin   = "100044238258"
  role_ids = [
    "308335260274237440",
    "308335260844662784"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String, ForceNew) Project ID.
* `role_ids` - (Required, Set: [`String`]) Role ID.
* `user_uin` - (Required, String, ForceNew) User ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateProjectMember](https://cloud.tencent.com/document/api/1267/123974) - Add project user role.
* Read: [ListProjectMembers](https://cloud.tencent.com/document/api/1267/123964) - This API is used to obtain the user under the project with pagination return.
* Update: [GrantMemberProjectRole](https://cloud.tencent.com/document/api/1267/123966) - Modify project user roles.
* Update: [RemoveMemberProjectRole](https://cloud.tencent.com/document/api/1267/123960) - Delete project user roles.
* Delete: [DeleteProjectMember](https://cloud.tencent.com/document/api/1267/123972) - This API is used to delete project users.

## Import

WeData project member can be imported using the projectId#userUin, e.g.

```
terraform import tencentcloud_wedata_project_member.example 2983848457986924544#100044238258
```

