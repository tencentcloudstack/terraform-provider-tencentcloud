Use this data source to query API gateway access keys.

Example Usage

```hcl
resource "tencentcloud_api_gateway_api_key" "test" {
  secret_name = "my_api_key"
  status      = "on"
}

data "tencentcloud_api_gateway_api_keys" "name" {
  secret_name = tencentcloud_api_gateway_api_key.test.secret_name
}

data "tencentcloud_api_gateway_api_keys" "id" {
  api_key_id = tencentcloud_api_gateway_api_key.test.id
}
```