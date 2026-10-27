Provides a resource to create a VDB instance.

~> **NOTE:** When `force_delete` is false (default), destroying this resource will only isolate the instance to the recycle bin. Set `force_delete` to true to permanently destroy the instance.

Example Usage

Create a pay-as-you-go single instance

```hcl
resource "tencentcloud_vdb_instance" "example" {
  vpc_id             = "vpc-i5yyodl9"
  subnet_id          = "subnet-hhi88a58"
  pay_mode           = 0
  security_group_ids = ["sg-4rd5741x"]
  instance_name      = "tf-example"
  instance_type      = "single"
  node_type          = "store"
  cpu                = 2
  memory             = 16
  disk_size          = 50
  worker_node_num    = 1
  force_delete       = true
}
```

Create a monthly subscription cluster instance with all parameters

```hcl
resource "tencentcloud_vdb_instance" "cluster" {
  vpc_id             = "vpc-i5yyodl9"
  subnet_id          = "subnet-hhi88a58"
  security_group_ids = ["sg-4rd5741x"]
  pay_mode           = 1
  pay_period         = 1
  auto_renew         = 1
  instance_name      = "tf-example"
  instance_type      = "cluster"
  mode               = "two"
  product_type       = 0
  node_type          = "compute"
  cpu                = 4
  memory             = 8
  disk_size          = 100
  worker_node_num    = 2
  params             = "{\"key\":\"value\"}"
  force_delete       = true
  resource_tags {
    tag_key   = "createBy"
    tag_value = "Terraform"
  }
}
```

Import

VDB instance can be imported using the id, e.g.

```
terraform import tencentcloud_vdb_instance.example vdb-oft69wcm
```
