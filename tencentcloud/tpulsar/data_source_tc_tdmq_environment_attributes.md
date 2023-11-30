Use this data source to query detailed information of tdmq environment_attributes

Example Usage

```hcl
data "tencentcloud_tdmq_environment_attributes" "example" {
  environment_id = tencentcloud_tdmq_namespace.example.environ_name
  cluster_id     = tencentcloud_tdmq_instance.example.id
}

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
  remark       = "remark."
}
```