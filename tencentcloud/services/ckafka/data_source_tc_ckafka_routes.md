Use this data source to query detailed information of CKafka routes

Example Usage

```hcl
data "tencentcloud_ckafka_routes" "example" {
  instance_id = "ckafka-exampleabc"
  main_route_flag = true
  route_id = 123
}

output "routes" {
  value = data.tencentcloud_ckafka_routes.example.routers
}
```
