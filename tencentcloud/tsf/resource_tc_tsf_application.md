Provides a resource to create a tsf application

Example Usage

```hcl
resource "tencentcloud_tsf_application" "application" {
  application_name = "my-app"
  application_type = "C"
  microservice_type = "M"
  application_desc = "This is my application"
  application_runtime_type = "Java"
  service_config_list {
		name = "my-service"
		ports {
			target_port = 8080
			protocol = "HTTP"
		}
		health_check {
			path = "/health"
		}
  }
  ignore_create_image_repository = true
}
```