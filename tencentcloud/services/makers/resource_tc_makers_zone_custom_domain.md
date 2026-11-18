subcategory: "Makers"
layout: "tencentcloudstack"
page_title: "TencentCloud: tencentcloud_makers_zone_custom_domain"
description: |-
  Provides a resource to bind a custom domain to a Makers (TEO-like pages service) zone.
---

# tencentcloud_makers_zone_custom_domain

Provides a resource to bind a custom domain to a Makers (TEO-like pages service) zone.

The domain is identified by the composite key `zone_id#project_id#domain`, which is
also the import ID format. Only the `env` attribute supports in-place update
(via `ModifyZoneCustomDomainEnv`); any other change forces replacement.

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

## Import

Makers zone custom domain can be imported using the id, e.g.

```sh
$ terraform import tencentcloud_makers_zone_custom_domain.example zone-123456#proj-abcdef#www.example.com
```
