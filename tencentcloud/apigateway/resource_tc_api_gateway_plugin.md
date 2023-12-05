Provides a resource to create a apiGateway plugin

Example Usage

```hcl
resource "tencentcloud_api_gateway_plugin" "example" {
  plugin_name = "tf-example"
  plugin_type = "IPControl"
  plugin_data = jsonencode({
    "type" : "white_list",
    "blocks" : "1.1.1.1",
  })
  description = "desc."
}
```

Import

apiGateway plugin can be imported using the id, e.g.

```
terraform import tencentcloud_api_gateway_plugin.plugin plugin_id
```