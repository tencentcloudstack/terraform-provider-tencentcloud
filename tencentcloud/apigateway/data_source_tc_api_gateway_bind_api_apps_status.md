Use this data source to query detailed information of apiGateway bind_api_apps_status

Example Usage

```hcl
data "tencentcloud_api_gateway_bind_api_apps_status" "example" {
  service_id = "service-nxz6yync"
  api_ids    = ["api-0cvmf4x4", "api-jvqlzolk"]
  filters {
    name   = "ApiAppId"
    values = ["app-krljp4wn"]
  }
}
```