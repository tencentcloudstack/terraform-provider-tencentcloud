---
subcategory: "Cloud Object Storage(COS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cos_bucket_version"
sidebar_current: "docs-tencentcloud-resource-cos_bucket_version"
description: |-
  Provides a resource to create a cos bucket_version
---

# tencentcloud_cos_bucket_version

Provides a resource to create a cos bucket_version

## Example Usage

```hcl
resource "tencentcloud_cos_bucket_version" "bucket_version" {
  bucket = "mycos-1258798060"
  status = "Enabled"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String) Bucket format should be [custom name]-[appid], for example `mycos-1258798060`.
* `status` - (Required, String) Whether to enable versioning. Valid values: `Suspended`, `Enabled`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cos bucket_version can be imported using the id, e.g.

```
terraform import tencentcloud_cos_bucket_version.bucket_version bucket_id
```

