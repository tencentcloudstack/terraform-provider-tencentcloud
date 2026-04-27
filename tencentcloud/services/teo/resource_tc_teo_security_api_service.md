Provides a resource to create a TEO security API service.

Example Usage

```hcl
resource "tencentcloud_teo_security_api_service" "example" {
  zone_id = "zone-2qtuhspy7cr6"

  api_services {
    name      = "my-api-service"
    base_path = "/api/v1"
  }

  api_resources {
    name            = "my-api-resource"
    api_service_ids = ["svc-id-12345"]
    path            = "/api/v1/users"
    methods         = ["GET", "POST"]
  }
}
```

Import

TEO security API service can be imported using the joint id "zone_id#api_service_ids", e.g.

```
terraform import tencentcloud_teo_security_api_service.example zone-2qtuhspy7cr6#svc-id-12345,svc-id-67890
```
