Provides a resource to create a scf reserved_concurrency_config

Example Usage

```hcl
resource "tencentcloud_scf_reserved_concurrency_config" "reserved_concurrency_config" {
  function_name = "keep-1676351130"
  reserved_concurrency_mem = 128000
  namespace     = "default"
}
```

Import

scf reserved_concurrency_config can be imported using the id, e.g.

```
terraform import tencentcloud_scf_reserved_concurrency_config.reserved_concurrency_config reserved_concurrency_config_id
```