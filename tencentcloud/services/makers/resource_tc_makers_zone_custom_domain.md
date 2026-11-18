Provides a resource to bind a custom domain to a Makers (TEO-like pages service) zone.

The domain is identified by the composite key `zone_id#project_id#domain`, which is
also the import ID format. Only the `env` attribute supports in-place update
(via `ModifyZoneCustomDomainEnv`); any other change forces replacement.

The binding, unbinding and environment switch operations are asynchronous on the
cloud side: the provider polls `DescribeZoneCustomDomains` until the change is
visible before the operation is reported as done. The polling budgets are
controlled by `create_timeout`, `update_timeout` and `delete_timeout`
(in minutes); when unset they default to the provider-level write timeout.

Note: the `DescribeZoneCustomDomains` API does not echo back the `project_id`,
so drift detection on `project_id` is limited to the value kept in state.

Example Usage

```hcl
resource "tencentcloud_makers_zone_custom_domain" "example" {
  zone_id        = "zone-123456"
  project_id     = "proj-abcdef"
  domain         = "www.example.com"
  env            = "production"
  create_timeout = 30
  update_timeout = 30
  delete_timeout = 30
}
```

Import

Makers zone custom domain can be imported using the id, e.g.

```sh
$ terraform import tencentcloud_makers_zone_custom_domain.example zone-123456#proj-abcdef#www.example.com
```
