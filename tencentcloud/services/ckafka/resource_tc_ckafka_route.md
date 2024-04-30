Provides a resource to create a ckafka route

Example Usage

```hcl
resource "tencentcloud_ckafka_route" "example" {
  instance_id    = "ckafka-8j4rodrr"
  vip_type       = 3
  vpc_id         = "vpc-axrsmmrv"
  subnet_id      = "subnet-j5vja918"
  access_type    = 0
  public_network = 3
}
```

Import

ckafka route can be imported using the id, e.g.

```
terraform import tencentcloud_ckafka_route.example ckafka-8j4rodrr#135912
```
