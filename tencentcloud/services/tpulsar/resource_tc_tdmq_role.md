Provide a resource to create a TDMQ role.

Example Usage

```hcl
resource "tencentcloud_tdmq_instance" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
  tags         = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tdmq_role" "example" {
  role_name  = "role_example"
  cluster_id = tencentcloud_tdmq_instance.example.id
  remark     = "remark."
}
```
