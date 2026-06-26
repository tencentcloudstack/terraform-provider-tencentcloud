Provides a resource to create a DBDC node to db custom cluster attachment.

Example Usage

```hcl
# create cluster
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

# create node
resource "tencentcloud_dbdc_db_custom_node" "example" {
  zone       = "ap-shanghai-5"
  image_id   = "img-rm13akp3"
  vpc_id     = "vpc-py7mlxqm"
  subnet_id  = "subnet-qd4upp83"
  node_type  = "DB.AT5.8XLARGE128"
  period     = 1
  auto_renew = 0
  node_name  = "tf-example"

  login_settings {
    password = "Password@2026"
  }

  tags = {
    createBy = "Terraform"
  }
}

# attach node to cluster
resource "tencentcloud_dbdc_node_to_db_custom_cluster_attachment" "example" {
  cluster_id = tencentcloud_dbdc_db_custom_cluster.example.id
  node_id    = tencentcloud_dbdc_db_custom_node.example.id
  image_id   = "img-rm13akp3"

  login_settings {
    password = "Passw0rd@2026"
  }
}
```

Import

DBDC node to db custom cluster attachment can be imported using the clusterId#nodeId, e.g.

```
terraform import tencentcloud_dbdc_node_to_db_custom_cluster_attachment.example dbcc-7uh7ludb#dbcn-ttiyh58n
```
