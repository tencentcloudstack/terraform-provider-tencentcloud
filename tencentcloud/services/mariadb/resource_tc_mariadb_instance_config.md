Provides a resource to create a mariadb instance_config

Example Usage

```hcl
resource "tencentcloud_mariadb_instance_config" "test" {
  instance_id        = "tdsql-9vqvls95"
  vpc_id             = "vpc-ii1jfbhl"
  subnet_id          = "subnet-3ku415by"
  rs_access_strategy = 1
  extranet_access    = 0
  vip                = "127.0.0.1"
}
```

Import

mariadb instance_config can be imported using the id, e.g.

```
terraform import tencentcloud_mariadb_instance_config.test id
```