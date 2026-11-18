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

Example Usage

```hcl
resource "tencentcloud_makers_zone_custom_domain" "example" {
  zone_id       = "zone-3hzk0ndvnnyp"
  project_id    = "makers-jw7jha6dktob"
  domain     = "tf.example.cn"
  env        = "Production"

  timeouts {
    create = "30m"
    update = "30m"
    delete = "30m"
  }
}
```

Import

Makers zone custom domain can be imported using the id, e.g.

```sh
terraform import tencentcloud_makers_zone_custom_domain.example zone-123456#proj-abcdef#www.example.com
```
