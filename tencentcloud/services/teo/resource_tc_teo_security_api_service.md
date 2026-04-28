Provides a resource to create a TEO API security service.

Example Usage

```hcl
resource "tencentcloud_teo_security_api_service" "example" {
  zone_id = "zone-3fkff38fyw8s"

  api_services {
    name      = "tf-example"
    base_path = "/api/v1"
  }
}
```

Import

TEO security API service can be imported using the zoneId#apiServiceId, e.g.

```
terraform import tencentcloud_teo_security_api_service.example zone-3fkff38fyw8s#apisrv-0000038524
```
