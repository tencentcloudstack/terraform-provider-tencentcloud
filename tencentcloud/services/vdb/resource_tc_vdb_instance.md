Provides a resource to create a VDB instance.

~> **NOTE:** When `force_delete` is false (default), destroying this resource will only isolate the instance to the recycle bin. Set `force_delete` to true to permanently destroy the instance.

Example Usage

Create a pay-as-you-go single instance

```hcl
resource "tencentcloud_vdb_instance" "example" {
  vpc_id             = "vpc-xxxxxxxx"
  subnet_id          = "subnet-xxxxxxxx"
  pay_mode           = 0
  security_group_ids = ["sg-xxxxxxxx"]
  instance_name      = "tf-example"
  instance_type      = "single"
  node_type          = "normal"
  cpu                = 2
  memory             = 4
  disk_size          = 50
  worker_node_num    = 1
  force_delete       = false
}
```

Create a monthly subscription cluster instance with all parameters

```hcl
resource "tencentcloud_vdb_instance" "cluster" {
  vpc_id            = "vpc-xxxxxxxx"
  subnet_id         = "subnet-xxxxxxxx"
  pay_mode          = 1
  pay_period        = 1
  auto_renew        = 1
  instance_name     = "tf-example-cluster"
  instance_type     = "cluster"
  mode              = "two"
  product_type      = 0
  node_type         = "compute"
  cpu               = 4
  memory            = 8
  disk_size         = 100
  worker_node_num   = 2
  params            = "{\"key\":\"value\"}"
  force_delete      = true

  security_group_ids = ["sg-xxxxxxxx"]

  resource_tags {
    tag_key   = "env"
    tag_value = "test"
  }

  resource_tags {
    tag_key   = "project"
    tag_value = "demo"
  }
}
```

Import

VDB instance can be imported using the id, e.g.

```
terraform import tencentcloud_vdb_instance.example vdb-xxxxxxxx
```
