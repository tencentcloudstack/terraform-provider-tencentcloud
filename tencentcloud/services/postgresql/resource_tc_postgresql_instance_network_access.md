Provides a resource to create a postgres instance network access

Example Usage

Create by custom vip

```hcl
resource "tencentcloud_postgresql_instance_network_access" "example" {
  db_instance_id = "postgres-ai46555b"
  vpc_id         = "vpc-i5yyodl9"
  subnet_id      = "subnet-d4umunpy"
  vip            = "10.0.10.11"
}
```

Create by automatic allocation vip

```hcl
resource "tencentcloud_postgresql_instance_network_access" "example" {
  db_instance_id = "postgres-ai46555b"
  vpc_id         = "vpc-i5yyodl9"
  subnet_id      = "subnet-d4umunpy"
}
```

Import

postgres instance network access can be imported using the id, e.g.

```
terraform import tencentcloud_postgresql_instance_network_access.example postgres-ai46555b#vpc-i5yyodl9#subnet-d4umunpy#10.0.10.11
```
