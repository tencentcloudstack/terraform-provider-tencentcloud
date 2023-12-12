Provides a resource to create a trocket rocketmq_role

Example Usage

```hcl
resource "tencentcloud_trocket_rocketmq_instance" "rocketmq_instance" {
  instance_type = "EXPERIMENT"
  name          = "test_role"
  sku_code      = "experiment_500"
  remark        = "test"
  vpc_id        = "vpc-xxxxx"
  subnet_id     = "subnet-xxxxx"
  tags          = {
    tag_key   = "rocketmq"
    tag_value = "5.x"
  }
}

resource "tencentcloud_trocket_rocketmq_role" "rocketmq_role" {
  instance_id = tencentcloud_trocket_rocketmq_instance.rocketmq_instance.id
  role        = "test_role"
  remark      = "test for terraform"
  perm_write  = false
  perm_read   = true
}

output "access_key" {
  value = tencentcloud_trocket_rocketmq_role.rocketmq_role.access_key
}

output "secret_key" {
  value = tencentcloud_trocket_rocketmq_role.rocketmq_role.secret_key
}
```

Import

trocket rocketmq_role can be imported using the id, e.g.

```
terraform import tencentcloud_trocket_rocketmq_role.rocketmq_role instanceId#role
```