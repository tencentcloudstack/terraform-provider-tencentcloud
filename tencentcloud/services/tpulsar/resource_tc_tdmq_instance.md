Provide a resource to create a TDMQ instance.

Example Usage

```hcl
resource "tencentcloud_tdmq_instance" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

Tdmq instance can be imported, e.g.

```
$ terraform import tencentcloud_tdmq_instance.example pulsar-78bwjaj8epxv
```