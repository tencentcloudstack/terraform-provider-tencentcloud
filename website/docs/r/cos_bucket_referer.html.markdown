---
subcategory: "Cloud Object Storage(COS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cos_bucket_referer"
sidebar_current: "docs-tencentcloud-resource-cos_bucket_referer"
description: |-
  Provides a resource to create a cos bucket_referer
---

# tencentcloud_cos_bucket_referer

Provides a resource to create a cos bucket_referer

## Example Usage

```hcl
resource "tencentcloud_cos_bucket_referer" "bucket_referer" {
  bucket                    = "mycos-1258798060"
  status                    = "Enabled"
  referer_type              = "Black-List"
  domain_list               = ["127.0.0.1", "terraform.com"]
  empty_refer_configuration = "Allow"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String) Bucket format should be [custom name]-[appid], for example `mycos-1258798060`.
* `domain_list` - (Required, Set: [`String`]) A list of domain names in the blocklist/allowlist.
* `referer_type` - (Required, String) Hotlink protection type. Enumerated values: `Black-List`, `White-List`.
* `status` - (Required, String) Whether to enable hotlink protection. Enumerated values: `Enabled`, `Disabled`.
* `empty_refer_configuration` - (Optional, String) Whether to allow access with an empty referer. Enumerated values: `Allow`, `Deny` (default).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cos bucket_referer can be imported using the id, e.g.

```
terraform import tencentcloud_cos_bucket_referer.bucket_referer bucket_id
```

