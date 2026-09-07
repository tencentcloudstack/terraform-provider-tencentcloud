Provides a resource to create a TDMQ professional cluster

Example Usage

Single-zone Professional Cluster

```hcl
resource "tencentcloud_tdmq_professional_cluster" "example" {
  auto_renew_flag  = 1
  cluster_name     = "tf_example"
  product_name     = "PULSAR.P2.SMALL4"
  storage_size     = 600
  instance_version = "2.9.2"
  zone_ids = [
    100006,
  ]

  vpc {
    vpc_id    = "vpc-i5yyodl9"
    subnet_id = "subnet-hhi88a58"
  }

  tags = {
    createby = "Terrafrom"
  }
}
```

Multi-zone Professional Cluster

```hcl
resource "tencentcloud_tdmq_professional_cluster" "example" {
  auto_renew_flag  = 1
  cluster_name     = "tf_example"
  product_name     = "PULSAR.P2.SMALL4"
  storage_size     = 600
  instance_version = "3.0.0"
  zone_ids = [
    100006,
    100007
  ]

  vpc {
    vpc_id    = "vpc-i5yyodl9"
    subnet_id = "subnet-hhi88a58"
  }

  tags = {
    createby = "Terrafrom"
  }
}
```

Import

TDMQ professional cluster can be imported using the id, e.g.

```
terraform import tencentcloud_tdmq_professional_cluster.example pulsar-x4r939zkwmm2
```
