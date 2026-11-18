---
subcategory: "Makers(makers)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_makers_zone_custom_domain"
sidebar_current: "docs-tencentcloud-resource-makers_zone_custom_domain"
description: |-
  Provides a resource to bind a custom domain to a Makers (TEO-like pages service) zone.
---

# tencentcloud_makers_zone_custom_domain

Provides a resource to bind a custom domain to a Makers (TEO-like pages service) zone.

The domain is identified by the composite key `zone_id#project_id#domain`, which is
also the import ID format. Only the `env` attribute supports in-place update
(via `ModifyZoneCustomDomainEnv`); any other change forces replacement.

The binding, unbinding and environment switch operations are asynchronous on the
cloud side: the provider polls `DescribeZoneCustomDomains` until the change is
visible before the operation is reported as done.

Note: the `DescribeZoneCustomDomains` API does not echo back the `project_id`,
so drift detection on `project_id` is limited to the value kept in state.

## Example Usage

```hcl
resource "tencentcloud_makers_zone_custom_domain" "example" {
  zone_id    = "zone-123456"
  project_id = "proj-abcdef"
  domain     = "www.example.com"
  env        = "production"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String) Custom domain name.
* `project_id` - (Required, String) Project ID.
* `zone_id` - (Required, String) Makers zone ID.
* `env` - (Optional, String) Environment slug the domain is bound to. If omitted on create, the backend decides the environment.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

      * `content` - Verification file content.
      * `name_servers` - Name servers.
      * `path` - Verification file path.
      * `record_type` - Record type, e.g. TXT.
      * `record_value` - Record value.
      * `subdomain` - Host record.
    * `dns_verification` - DNS verification.
    * `file_verification` - File verification.
    * `include_sub_domains` - Whether to include sub domains. Valid values: `on`, `off`.
    * `max_age` - Validity period in seconds.
    * `ns_verification` - NS verification.
    * `preload` - Whether to preload. Valid values: `on`, `off`.
    * `redirect_status_code` - Redirect status code. Valid values: 301, 302, 307, 308.
    * `switch` - Switch state. Valid values: `on`, `off`.
    * `switch` - Switch state. Valid values: `on`, `off`.
    * `switch` - Switch state. Valid values: `on`, `off`.
    * `switch` - Switch state. Valid values: `on`, `off`.
    * `switch` - Switch state. Valid values: `on`, `off`.
  * `area` - Acceleration area.
  * `cname` - CNAME.
  * `domain` - Domain name.
  * `env` - Bound environment slug, defaults to Production.
  * `force_redirect` - Force HTTPS redirect configuration.
  * `hsts` - HSTS configuration.
  * `http2` - HTTP/2 configuration.
  * `ipv6` - IPv6 configuration.
  * `main_domain` - Main domain.
  * `ocsp_stapling` - OCSP stapling configuration.
  * `ownership_verification` - Ownership verification information.
  * `status` - Domain status.
  * `type` - Domain type. Valid values: `Custom`, `Preset`.
  * `zone_id` - Zone ID.
* `id` - ID of the resource in the format `zone_id#project_id#domain`.
* `pages_domains` - The current domain configuration returned by the cloud side. Only the managed custom domain is reported here.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateZoneCustomDomain](https://www.tencentcloud.com/search/?wd=CreateZoneCustomDomain&type=document&lang=en)
  - [DescribeZoneCustomDomains](https://www.tencentcloud.com/search/?wd=DescribeZoneCustomDomains&type=document&lang=en)
* Read: [DescribeZoneCustomDomains](https://www.tencentcloud.com/search/?wd=DescribeZoneCustomDomains&type=document&lang=en)
* Update:
  - [DescribeZoneCustomDomains](https://www.tencentcloud.com/search/?wd=DescribeZoneCustomDomains&type=document&lang=en)
  - [ModifyZoneCustomDomainEnv](https://www.tencentcloud.com/search/?wd=ModifyZoneCustomDomainEnv&type=document&lang=en)
* Delete:
  - [DeleteProjectCustomDomains](https://www.tencentcloud.com/search/?wd=DeleteProjectCustomDomains&type=document&lang=en)
  - [DescribeZoneCustomDomains](https://www.tencentcloud.com/search/?wd=DescribeZoneCustomDomains&type=document&lang=en)

## Import

Makers zone custom domain can be imported using the id, e.g.

```sh
$ terraform import tencentcloud_makers_zone_custom_domain.example zone-123456#proj-abcdef#www.example.com
```

