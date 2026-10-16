---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_alias_domain"
sidebar_current: "docs-tencentcloud-resource-teo_alias_domain"
description: |-
  Provides a resource to create a TEO alias domain.
---

# tencentcloud_teo_alias_domain

Provides a resource to create a TEO alias domain.

~> **NOTE:** This feature is only supported by the Enterprise edition plan and is currently in beta testing. Please [contact us](https://cloud.tencent.com/online-service?from=connect-us) if you need to use it.

## Example Usage

```hcl
resource "tencentcloud_teo_alias_domain" "example" {
  zone_id     = "zone-3fkff38fyw8s"
  alias_name  = "alias.demo.com"
  target_name = "target.demo.com"
  cert_type   = "none"
}
```

### With SSL hosted certificate

```hcl
resource "tencentcloud_teo_alias_domain" "example" {
  zone_id     = "zone-3fkff38fyw8s"
  alias_name  = "alias.demo.com"
  target_name = "target.demo.com"
  cert_type   = "hosting"
  cert_id     = ["your-cert-id"]
}
```

### With disabled status

```hcl
resource "tencentcloud_teo_alias_domain" "example" {
  zone_id     = "zone-3fkff38fyw8s"
  alias_name  = "alias.demo.com"
  target_name = "target.demo.com"
  cert_type   = "none"
  paused      = true
}
```

## Argument Reference

The following arguments are supported:

* `alias_name` - (Required, String, ForceNew) Alias domain name.
* `target_name` - (Required, String) Target domain name.
* `zone_id` - (Required, String, ForceNew) Site ID.
* `cert_id` - (Optional, List: [`String`]) Certificate ID list. Required when `cert_type` is `hosting`.
* `cert_type` - (Optional, String) Certificate configuration. Valid values: `none` (no configuration), `hosting` (SSL hosted certificate). Default value: `none`.
* `paused` - (Optional, Bool) Whether the alias domain is disabled. `false`: enabled; `true`: disabled.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_on` - Alias domain creation time.
* `forbid_mode` - Block mode. Valid values: `0` (not blocked), `11` (compliance blocked), `14` (not registered blocked).
* `modified_on` - Alias domain modification time.
* `status` - Alias domain status. Valid values: `active` (effective), `pending` (deploying), `conflict` (reclaimed), `stop` (disabled).

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `10m`) Used when creating the resource.
* `update` - (Defaults to `10m`) Used when updating the resource.
## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateAliasDomain](https://www.tencentcloud.com/search/?wd=CreateAliasDomain&type=document&lang=en)
* Create: [DescribeAliasDomains](https://www.tencentcloud.com/search/?wd=DescribeAliasDomains&type=document&lang=en)
* Create: [ModifyAliasDomainStatus](https://www.tencentcloud.com/search/?wd=ModifyAliasDomainStatus&type=document&lang=en)
* Read: [DescribeAliasDomains](https://www.tencentcloud.com/search/?wd=DescribeAliasDomains&type=document&lang=en)
* Update: [DescribeAliasDomains](https://www.tencentcloud.com/search/?wd=DescribeAliasDomains&type=document&lang=en)
* Update: [ModifyAliasDomain](https://www.tencentcloud.com/search/?wd=ModifyAliasDomain&type=document&lang=en)
* Update: [ModifyAliasDomainStatus](https://www.tencentcloud.com/search/?wd=ModifyAliasDomainStatus&type=document&lang=en)
* Delete: [DeleteAliasDomain](https://www.tencentcloud.com/search/?wd=DeleteAliasDomain&type=document&lang=en)
* Delete: [DescribeAliasDomains](https://www.tencentcloud.com/search/?wd=DescribeAliasDomains&type=document&lang=en)
* Delete: [ModifyAliasDomainStatus](https://www.tencentcloud.com/search/?wd=ModifyAliasDomainStatus&type=document&lang=en)

