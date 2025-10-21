---
subcategory: "Cloud Infinite(CI)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ci_bucket_attachment"
sidebar_current: "docs-tencentcloud-resource-ci_bucket_attachment"
description: |-
  Provides a resource to create a ci bucket
---

# tencentcloud_ci_bucket_attachment

Provides a resource to create a ci bucket

## Example Usage

```hcl
resource "tencentcloud_ci_bucket_attachment" "bucket_attachment" {
  bucket = "terraform-ci-xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String, ForceNew) bucket name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `ci_status` - Binding object storage state, `on`: bound, `off`: unbound, `unbinding`: unbinding.


## Import

ci bucket can be imported using the id, e.g.

```
terraform import tencentcloud_ci_bucket_attachment.bucket_attachment terraform-ci-xxxxxx
```

