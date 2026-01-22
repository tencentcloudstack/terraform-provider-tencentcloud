Provide a resource to create a Mongodb instance.

Example Usage

```hcl
resource "tencentcloud_mongodb_instance" "example" {
  instance_name  = "tf-example"
  memory         = 4
  volume         = 100
  engine_version = "MONGO_40_WT"
  machine_type   = "HIO10G"
  available_zone = "ap-guangzhou-6"
  vpc_id         = "vpc-i5yyodl9"
  subnet_id      = "subnet-hhi88a58"
  project_id     = 0
  password       = "Password@123"
}
```

Import

Mongodb instance can be imported using the id, e.g.

```
terraform import tencentcloud_mongodb_instance.example cmgo-41s6jwy4
```