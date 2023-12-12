Provides a resource to create a apiGateway import_open_api

Example Usage

Import open Api by YAML

```hcl
resource "tencentcloud_api_gateway_import_open_api" "example" {
  service_id      = "service-nxz6yync"
  content         = "info:\n  title: keep-service\n  version: 1.0.1\nopenapi: 3.0.0\npaths:\n  /api/test:\n    get:\n      description: desc\n      operationId: test\n      responses:\n        '200':\n          content:\n            text/html:\n              example: '200'\n          description: '200'\n        default:\n          content:\n            text/html:\n              example: '400'\n          description: '400'\n      x-apigw-api-business-type: NORMAL\n      x-apigw-api-type: NORMAL\n      x-apigw-backend:\n        ServiceConfig:\n          Method: GET\n          Path: /test\n          Url: http://domain.com\n        ServiceType: HTTP\n      x-apigw-cors: false\n      x-apigw-protocol: HTTP\n      x-apigw-service-timeout: 15\n"
  encode_type     = "YAML"
  content_version = "openAPI"
}
```

Import open Api by JSON

```hcl
resource "tencentcloud_api_gateway_import_open_api" "example" {
  service_id      = "service-nxz6yync"
  content         = "{\"openapi\": \"3.0.0\", \"info\": {\"title\": \"keep-service\", \"version\": \"1.0.1\"}, \"paths\": {\"/api/test\": {\"get\": {\"operationId\": \"test\", \"description\": \"desc\", \"responses\": {\"200\": {\"description\": \"200\", \"content\": {\"text/html\": {\"example\": \"200\"}}}, \"default\": {\"content\": {\"text/html\": {\"example\": \"400\"}}, \"description\": \"400\"}}, \"x-apigw-api-type\": \"NORMAL\", \"x-apigw-api-business-type\": \"NORMAL\", \"x-apigw-protocol\": \"HTTP\", \"x-apigw-cors\": false, \"x-apigw-service-timeout\": 15, \"x-apigw-backend\": {\"ServiceType\": \"HTTP\", \"ServiceConfig\": {\"Url\": \"http://domain.com\", \"Path\": \"/test\", \"Method\": \"GET\"}}}}}}"
  encode_type     = "JSON"
  content_version = "openAPI"
}
```