Provides a resource to manage a cloud native API gateway AI service source.

Example Usage

```hcl
resource "tencentcloud_cngw_cloud_native_api_gateway_ai_service_source" "example" {
  gateway_id  = "gw-12345"
  source_type = "MCPRegistry"
  source_name = "tf-example-source"
  source_info = {
    instance_id = "ins-12345"
    addresses   = ["10.0.0.1:8848"]
  }
  description = "created by terraform"
}
```

Import

Cloud native API gateway AI service source can be imported using the id, e.g.

```
terraform import tencentcloud_cngw_cloud_native_api_gateway_ai_service_source.example source-12345
```
