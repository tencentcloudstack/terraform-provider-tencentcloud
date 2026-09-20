Provides a resource to create a TEO inference API token.

Example Usage

```hcl
resource "tencentcloud_teo_inference_api_token" "example" {
  zone_id = "zone-2qtuhspy7cr6"
  name    = "my-inference-token"
}
```

Import

TEO inference API token can be imported using the joint id "zone_id#token_id", e.g.

```
terraform import tencentcloud_teo_inference_api_token.example zone-2qtuhspy7cr6#token-abcdefghij
```
