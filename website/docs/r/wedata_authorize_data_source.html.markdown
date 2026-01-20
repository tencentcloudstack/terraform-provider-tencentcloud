---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_authorize_data_source"
sidebar_current: "docs-tencentcloud-resource-wedata_authorize_data_source"
description: |-
  Provides a resource to create a WeData authorize data source
---

# tencentcloud_wedata_authorize_data_source

Provides a resource to create a WeData authorize data source

## Example Usage

### Authorize by project ids

```hcl
resource "tencentcloud_wedata_authorize_data_source" "example" {
  data_source_id = "116203"
  auth_project_ids = [
    "1857740139240632320",
    "1857740139240632318",
  ]
}
```

### Authorize by users

```hcl
resource "tencentcloud_wedata_authorize_data_source" "example" {
  data_source_id = "116203"
  auth_users = [
    "1857740139240632320_100028448903",
    "1857740139240632320_100028578751",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `data_source_id` - (Required, String, ForceNew) Data source ID.
* `auth_project_ids` - (Optional, Set: [`String`]) List of target project ID to be authorized.
* `auth_users` - (Optional, Set: [`String`]) List of users under the authorized project, format: project_id_user_id.
When authorizing multiple objects, the project ID must be consistent.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

WeData authorize data source can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_authorize_data_source.example 116203
```

