Provides a resource to create a scf function_event_invoke_config

Example Usage

```hcl
resource "tencentcloud_scf_function_event_invoke_config" "function_event_invoke_config" {
  function_name = "keep-1676351130"
  namespace     = "default"
  async_trigger_config {
    retry_config {
      retry_num = 2
    }
    msg_ttl = 24
  }
}
```

Import

scf function_event_invoke_config can be imported using the id, e.g.

```
terraform import tencentcloud_scf_function_event_invoke_config.function_event_invoke_config function_name#namespace
```