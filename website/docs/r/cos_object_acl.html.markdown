---
subcategory: "Cloud Object Storage(COS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cos_object_acl"
sidebar_current: "docs-tencentcloud-resource-cos_object_acl"
description: |-
  Provides a resource to create a cos object_acl
---

# tencentcloud_cos_object_acl

Provides a resource to create a cos object_acl

## Example Usage

```hcl
resource "tencentcloud_cos_object_acl" "object_acl" {
  bucket    = "mycos-1258798060"
  key       = "exampleobject"
  x_cos_acl = "private"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String) Bucket format should be [custom name]-[appid], for example `mycos-1258798060`.
* `key` - (Required, String) object name in bucket.
* `x_cos_acl` - (Required, String) ACL attribute of the object.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cos object_acl can be imported using the id, e.g.

```
terraform import tencentcloud_cos_object_acl.object_acl object_acl_id
```

