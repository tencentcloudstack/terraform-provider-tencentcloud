---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_ownership_verify"
sidebar_current: "docs-tencentcloud-resource-teo_ownership_verify"
description: |-
  Provides a resource to create a teo ownership_verify
---

# tencentcloud_teo_ownership_verify

Provides a resource to create a teo ownership_verify

## Example Usage

```hcl
resource "tencentcloud_teo_ownership_verify" "ownership_verify" {
  domain = "qq.com"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) Verify domain name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `result` - When the verification result is failed, this field will return the reason.
* `status` - Ownership verification results. `success`: verification successful; `fail`: verification failed.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [VerifyOwnership](https://cloud.tencent.com/document/api/1552/98879) - This API is used to verify your ownership of a site or domain name. It's required in the CNAME access mode. After a site is verified, you don't need to verify the ownership again for domain names added to it in the future. For details, see [Ownership Verification](https://intl.cloud.tencent.com/document/product/1552/70789?from_cn_redirect=1).

