Provides a resource to create a DBDC db custom node.

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

resource "tencentcloud_dbdc_db_custom_node" "example" {
  zone       = "ap-shanghai-5"
  image_id   = "img-rm13akp3"
  vpc_id     = "vpc-py7mlxqm"
  subnet_id  = "subnet-qd4upp83"
  node_type  = "DB.AT5.8XLARGE128"
  period     = 1
  auto_renew = 2
  node_name  = "tf-example"

  login_settings {
    password = "Password@2026"
  }

  tags = {
    createBy = "Terraform"
  }
}
```

Import

DBDC db custom node can be imported using the id, e.g.

```
terraform import tencentcloud_dbdc_db_custom_node.example dbcn-ttiyh58n
```
