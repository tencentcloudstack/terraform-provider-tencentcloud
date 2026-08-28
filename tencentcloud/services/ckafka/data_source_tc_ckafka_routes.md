Use this data source to query detailed information of CKafka routes

Example Usage

```hcl
data "tencentcloud_ckafka_routes" "example" {
  instance_id = "ckafka-bqwlyrg8"
}

output "routes" {
  value = data.tencentcloud_ckafka_routes.example.routers
}
```
