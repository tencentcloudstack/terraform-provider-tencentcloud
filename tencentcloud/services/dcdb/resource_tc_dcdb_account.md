Provides a resource to create a dcdb account

Example Usage

```hcl
resource "tencentcloud_dcdb_hourdb_instance" "example" {
  instance_name     = "tf-example"
  zones             = ["ap-guangzhou-6", "ap-guangzhou-7"]
  shard_memory      = "4"
  shard_storage     = "100"
  shard_node_count  = "2"
  shard_count       = "2"
  vpc_id            = "vpc-i5yyodl9"
  subnet_id         = "subnet-hhi88a58"
  security_group_id = "sg-4z20n68d"
  db_version_id     = "8.0"
  resource_tags {
    tag_key   = "tagKey"
    tag_value = "tagValue"
  }
}

resource "tencentcloud_dcdb_account" "example" {
  instance_id          = tencentcloud_dcdb_hourdb_instance.example.id
  user_name            = "mysql"
  host                 = "127.0.0.1"
  password             = "Password@2026"
  read_only            = 0
  description          = "this is a test account"
  max_user_connections = 10
}
```
Import

dcdb account can be imported using the instanceId#userName, e.g.
```
terraform import tencentcloud_dcdb_account.example tdsqlshard-0tyhkwkzj5#mysql
```
