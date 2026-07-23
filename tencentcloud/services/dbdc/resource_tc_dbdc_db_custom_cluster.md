Provides a resource to create a DBDC db custom cluster.

Example Usage

```hcl
resource "tencentcloud_dbdc_db_custom_cluster" "example" {
  cluster_name        = "tf-example"
  cluster_description = "cluster description."

  container_network {
    vpc_id     = "vpc-py7mlxqm"
    subnet_ids = ["subnet-qd4upp83", "subnet-g7vmz57f", "subnet-hqbm5bwx"]
  }

  api_server_network {
    vpc_id    = "vpc-b4zgfr3a"
    subnet_id = "subnet-cp3juq8r"
  }

  tags = {
    createBy = "Terraform"
  }
}
```

Import

DBDC db custom cluster can be imported using the id, e.g.

```
terraform import tencentcloud_dbdc_db_custom_cluster.example dbcc-k7snlxyu
```
