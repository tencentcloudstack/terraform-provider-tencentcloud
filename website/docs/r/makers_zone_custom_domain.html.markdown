---
subcategory: "Web & Agents Development Platform(Makers)"
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
visible before the operation is reported as done. The waiting budgets are
controlled through the standard `timeouts` block (`create`, `update` and `delete` 
attributes, duration strings such as `30m`); when unset they fall
back to the provider-level timeouts.

## Example Usage

```hcl
resource "tencentcloud_makers_zone_custom_domain" "example" {
  zone_id    = "zone-3hzk0ndvnnyp"
  project_id = "makers-jw7jha6dktob"
  domain     = "tf.example.cn"
  env        = "Production"

  timeouts {
    create = "30m"
    update = "30m"
    delete = "30m"
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String) Custom domain name.
* `project_id` - (Required, String) Project ID.
* `zone_id` - (Required, String) Makers zone ID.
* `env` - (Optional, String) Environment slug the domain is bound to. If omitted on create, the backend decides the environment. Production environment: `Production`; Preview: `Preview`.
* `timeouts` - (Optional, Object) The `timeouts` block.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource in the format `zone_id#project_id#domain`.

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
terraform import tencentcloud_makers_zone_custom_domain.example zone-123456#proj-abcdef#www.example.com
```

