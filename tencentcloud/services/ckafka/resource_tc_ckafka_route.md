Provides a resource to create a ckafka route

Example Usage

```hcl
resource "tencentcloud_ckafka_route" "example" {
  instance_id    = "ckafka-8j4ra928"
  vip_type       = 3
  vpc_id         = "vpc-i5yyodl9"
  subnet_id      = "subnet-hhi88a58"
  access_type    = 0
  public_network = 3
  security_group_ids = [
    "sg-8ztqghor",
    "sg-rc9644tr",
    "sg-n8zf5ry9"
  ]
}
```

Import

ckafka route can be imported using the instanceId#routeId, e.g.

```
terraform import tencentcloud_ckafka_route.example ckafka-8j4ra928#135912
```
