Provides a resource to create a tsf microservice

Example Usage

```hcl
resource "tencentcloud_tsf_microservice" "microservice" {
  namespace_id = "namespace-vjlkzkgy"
  microservice_name = "test-microservice"
  microservice_desc = "desc-microservice"
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

tsf microservice can be imported using the namespaceId#microserviceId, e.g.

```
terraform import tencentcloud_tsf_microservice.microservice namespace-vjlkzkgy#ms-vjeb43lw
```