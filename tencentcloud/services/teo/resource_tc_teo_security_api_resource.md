Provides a resource to create a TEO API security resource.

Example Usage

```hcl
resource "tencentcloud_teo_security_api_service" "example" {
  zone_id = "zone-3fkff38fyw8s"

  api_services {
    name      = "tf-example"
    base_path = "/api/v1"
  }
}

resource "tencentcloud_teo_security_api_resource" "example" {
  zone_id = "zone-3fkff38fyw8s"

  api_resources {
    name               = "tf-example"
    path               = "/api/v1/orders"
    api_service_ids    = [tencentcloud_teo_security_api_service.example.api_services[0].id]
    methods            = ["GET", "POST"]
    request_constraint = "$${http.request.body.form['operationType']} in ['query', 'create']"
  }
}
```

Import

TEO security API resource can be imported using the zoneId#apiResourceId, e.g.

```
terraform import tencentcloud_teo_security_api_resource.example zone-3fkff38fyw8s#apires-0000039154
```
