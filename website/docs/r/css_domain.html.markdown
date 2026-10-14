---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_domain"
sidebar_current: "docs-tencentcloud-resource-css_domain"
description: |-
  Provides a resource to create a css domain
---

# tencentcloud_css_domain

Provides a resource to create a css domain

## Example Usage

```hcl
resource "tencentcloud_css_domain" "domain" {
  domain_name          = "iac-tf.cloud"
  domain_type          = 0
  play_type            = 1
  is_delay_live        = 0
  is_mini_program_live = 0
  verify_owner_type    = "dbCheck"
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, String) Domain Name.
* `domain_type` - (Required, Int) Domain type: `0`: push stream. `1`: playback.
* `enable` - (Optional, Bool) Switch. true: enable the specified domain, false: disable the specified domain.
* `is_delay_live` - (Optional, Int) Whether it is LCB: `0`: LVB. `1`: LCB. Default value is 0.
* `is_mini_program_live` - (Optional, Int) `0`: LVB. `1`: LVB on Mini Program. Note: this field may return null, indicating that no valid values can be obtained. Default value is 0.
* `play_type` - (Optional, Int) Play Type. This parameter is valid only if `DomainType` is 1. Available values: `1`: in Mainland China. `2`: global. `3`: outside Mainland China. Default value is 1.
* `verify_owner_type` - (Optional, String) Domain name attribution verification type. `dnsCheck`, `fileCheck`, `dbCheck`. The default is `dbCheck`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [AddLiveDomain](https://www.tencentcloud.com/document/product/267/35189) - This API is used to add a domain name. Only one domain name can be submitted at a time, and it must have an ICP filing.
* Create: [AuthenticateDomainOwner](https://www.tencentcloud.com/document/product/267/50612) - This API is used to verify the ownership of a domain.
* Create: [EnableLiveDomain](https://www.tencentcloud.com/document/product/267/35185) - This API is used to enable a disabled LVB domain name.
* Create: [ForbidLiveDomain](https://www.tencentcloud.com/document/product/267/35184) - This API is used to disable an LVB domain name.
* Read: [DescribeLiveDomain](https://www.tencentcloud.com/document/product/267/35187) - This API is used to query the LVB domain name information.
* Update: [EnableLiveDomain](https://www.tencentcloud.com/document/product/267/35185) - This API is used to enable a disabled LVB domain name.
* Update: [ForbidLiveDomain](https://www.tencentcloud.com/document/product/267/35184) - This API is used to disable an LVB domain name.
* Update: [ModifyLivePlayDomain](https://www.tencentcloud.com/document/product/267/35183) - This API is used to modify a playback domain name.
* Delete: [DeleteLiveDomain](https://www.tencentcloud.com/document/product/267/35188) - This API is used to delete an added LVB domain name.

## Import

css domain can be imported using the id, e.g.

```
terraform import tencentcloud_css_domain.domain domain_name
```

