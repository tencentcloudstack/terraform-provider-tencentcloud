Use this data source to query detailed information of apiGateway api_plugins

Example Usage

```hcl
data "tencentcloud_api_gateway_api_plugins" "example" {
  api_id           = "api-0cvmf4x4"
  service_id       = "service-nxz6yync"
  environment_name = "test"
}
```