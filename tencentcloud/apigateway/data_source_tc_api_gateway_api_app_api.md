Use this data source to query detailed information of apiGateway api_app_api

Example Usage

```hcl
data "tencentcloud_api_gateway_api_app_api" "example" {
  service_id = "service-nxz6yync"
  api_id     = "api-0cvmf4x4"
  api_region = "ap-guangzhou"
}
```