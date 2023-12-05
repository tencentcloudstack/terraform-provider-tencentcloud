Provide a resource to create a Mongodb instance.

Example Usage

```hcl
resource "tencentcloud_mongodb_instance" "mongodb" {
  instance_name  = "mongodb"
  memory         = 4
  volume         = 100
  engine_version = "MONGO_36_WT"
  machine_type   = "HIO10G"
  available_zone = "ap-guangzhou-2"
  vpc_id         = "vpc-xxxxxx"
  subnet_id      = "subnet-xxxxxx"
  project_id     = 0
  password       = "password1234"
}
```

Import

Mongodb instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_mongodb_instance.mongodb cmgo-41s6jwy4
```