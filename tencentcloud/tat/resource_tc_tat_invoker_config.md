Provides a resource to create a tat invoker_config

Example Usage

```hcl
resource "tencentcloud_tat_invoker_config" "invoker_config" {
  invoker_id = "ivk-cas4upyf"
  invoker_status = "on"
}
```

Import

tat invoker_config can be imported using the id, e.g.

```
terraform import tencentcloud_tat_invoker_config.invoker_config invoker_config_id
```