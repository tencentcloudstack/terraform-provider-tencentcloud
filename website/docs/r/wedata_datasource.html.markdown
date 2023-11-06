---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_datasource"
sidebar_current: "docs-tencentcloud-resource-wedata_datasource"
description: |-
  Provides a resource to create a wedata datasource
---

# tencentcloud_wedata_datasource

Provides a resource to create a wedata datasource

## Example Usage

```hcl
resource "tencentcloud_wedata_datasource" "example" {
  name                = "tf_example"
  category            = "DB"
  type                = "MYSQL"
  owner_project_id    = "110111121"
  owner_project_name  = "ownerprojectname"
  owner_project_ident = "OwnerProjectIdent"
  biz_params          = "{}"
  params              = "{}"
  description         = "descr"
  display             = "Display"
  database_name       = "db"
  instance            = "instance"
  status              = 1
  cluster_id          = "cid"
  collect             = "false"
  cos_bucket          = "aaaa"
  cos_region          = "ap-guangzhou"
}
```

## Argument Reference

The following arguments are supported:

* `category` - (Required, String) DataSource Category.
* `name` - (Required, String) DataSource Name.
* `owner_project_id` - (Required, String) Owner projectId.
* `owner_project_ident` - (Required, String) Owner Project Ident.
* `owner_project_name` - (Required, String) Owner project name.
* `type` - (Required, String) DataSource Type.
* `biz_params` - (Optional, String) BizParams.
* `cluster_id` - (Optional, String) ClusterId.
* `collect` - (Optional, String) Collect.
* `cos_bucket` - (Optional, String) COSBucket.
* `cos_region` - (Optional, String) Cos region.
* `database_name` - (Optional, String) Dbname.
* `description` - (Optional, String) Description.
* `display` - (Optional, String) Display.
* `instance` - (Optional, String) Instance.
* `params` - (Optional, String) Params.
* `status` - (Optional, Int) Status.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



