Provides a resource to create a tem application_service

Example Usage

```hcl
resource "tencentcloud_tem_application_service" "application_service" {
  environment_id = "en-dpxyydl5"
  application_id = "app-jrl3346j"
  service {
		type = "CLUSTER"
		service_name = "test0-1"
		port_mapping_item_list {
			port = 80
			target_port = 80
			protocol = "TCP"
		}
  }
}
```

Import

tem application_service can be imported using the environmentId#applicationId#serviceName, e.g.

```
terraform import tencentcloud_tem_application_service.application_service en-dpxyydl5#app-jrl3346j#test0-1
```