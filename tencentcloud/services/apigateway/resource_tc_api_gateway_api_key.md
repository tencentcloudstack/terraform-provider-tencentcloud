Use this resource to create API gateway access key.

Example Usage

Automatically generate key for API gateway access key.

```hcl
resource "tencentcloud_api_gateway_api_key" "example_auto" {
  secret_name = "tf_example_auto"
  status      = "on"
}
```

Manually generate a secret key for API gateway access key.

```hcl
resource "tencentcloud_api_gateway_api_key" "example_manual" {
  secret_name       = "tf_example_manual"
  status            = "on"
  access_key_type   = "manual"
  access_key_id     = "28e287e340507fa147b2c8284dab542f"
  access_key_secret = "0198a4b8c3105080f4acd9e507599eff"
}
```
Import

API gateway access key can be imported using the id, e.g.

```
$ terraform import tencentcloud_api_gateway_api_key.test AKIDMZwceezso9ps5p8jkro8a9fwe1e7nzF2k50B
```