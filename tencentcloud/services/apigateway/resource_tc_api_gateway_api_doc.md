Provides a resource to create a APIGateway ApiDoc

Example Usage

```hcl
resource "tencentcloud_api_gateway_api_doc" "my_api_doc" {
  api_doc_name = "doc_test1"
  service_id   = "service_test1"
  environment  = "release"
  api_ids      = ["api-test1", "api-test2"]
}
```