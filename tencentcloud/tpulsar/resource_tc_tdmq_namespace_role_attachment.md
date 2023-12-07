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

resource "tencentcloud_tdmq_namespace" "example" {
  environ_name = "tf_example"
  msg_ttl      = 300
  cluster_id   = tencentcloud_tdmq_instance.example.id
  retention_policy {
    time_in_minutes = 60
    size_in_mb      = 10
  }
  remark = "remark."
}

resource "tencentcloud_tdmq_role" "example" {
  role_name  = "tf_example"
  cluster_id = tencentcloud_tdmq_instance.example.id
  remark     = "remark."
}

resource "tencentcloud_tdmq_namespace_role_attachment" "example" {
  environ_id  = tencentcloud_tdmq_namespace.example.environ_name
  role_name   = tencentcloud_tdmq_role.example.role_name
  permissions = ["produce", "consume"]
  cluster_id  = tencentcloud_tdmq_instance.example.id
}
```