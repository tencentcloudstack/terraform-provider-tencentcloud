Provide a resource to create a TDMQ instance.

Example Usage

```hcl
resource "tencentcloud_tdmq_instance" "foo" {
  cluster_name = "example"
  remark = "this is description."
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

Tdmq instance can be imported, e.g.

```
$ terraform import tencentcloud_tdmq_instance.test tdmq_id
```