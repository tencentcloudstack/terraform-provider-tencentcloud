---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_repository"
sidebar_current: "docs-tencentcloud-resource-tsf_repository"
description: |-
  Provides a resource to create a tsf repository
---

# tencentcloud_tsf_repository

Provides a resource to create a tsf repository

## Example Usage

```hcl
resource "tencentcloud_tsf_repository" "repository" {
  repository_name = ""
  repository_type = ""
  bucket_name     = ""
  bucket_region   = ""
  directory       = ""
  repository_desc = ""
}
```

## Argument Reference

The following arguments are supported:

* `bucket_name` - (Required, String) the name of the bucket where the warehouse is located.
* `bucket_region` - (Required, String) Bucket region where the warehouse is located.
* `repository_name` - (Required, String) warehouse name.
* `repository_type` - (Required, String) warehouse type (default warehouse: default, private warehouse: private).
* `directory` - (Optional, String) directory.
* `repository_desc` - (Optional, String) warehouse description.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - warehouse creation time.
* `is_used` - whether the repository is in use.
* `repository_id` - Warehouse ID.


## Import

tsf repository can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_repository.repository repository_id
```

