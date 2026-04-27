Provides a resource to create a TEO (EdgeOne) security API resource, which is used to define API endpoints and their associated API services for security protection.

Example Usage

```hcl
resource "tencentcloud_teo_security_api_resource" "example" {
  zone_id = "zone-2qtuhspy7cr6"
  api_resources {
    name             = "test-api-resource"
    api_service_ids  = ["svc-123"]
    path             = "/api/v1/test"
    methods          = ["GET", "POST"]
    request_constraint = jsonencode({"key": "value"})
  }
  api_resources {
    name    = "test-api-resource-2"
    path    = "/api/v2/test"
    methods = ["GET"]
  }
}
```

Import

teo security_api_resource can be imported using the id, e.g.

```
terraform import tencentcloud_teo_security_api_resource.example zone_id
```
