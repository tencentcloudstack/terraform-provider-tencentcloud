Provides a resource to create a APIGateway ApiApp

Example Usage

Create a basic apigateway api_app

```hcl
resource "tencentcloud_api_gateway_api_app" "example" {
  api_app_name = "tf_example"
  api_app_desc = "app desc."
}
```

Bind Tag

```hcl
resource "tencentcloud_api_gateway_api_app" "example" {
  api_app_name = "tf_example"
  api_app_desc = "app desc."

  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

apigateway api_app can be imported using the id, e.g.

```
terraform import tencentcloud_api_gateway_api_app.example app-poe0pyex
```