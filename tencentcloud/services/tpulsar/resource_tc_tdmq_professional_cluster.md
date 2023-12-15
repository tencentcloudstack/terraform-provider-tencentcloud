Provides a resource to create a tdmq professional_cluster

Example Usage

single-zone Professional Cluster

```hcl
resource "tencentcloud_tdmq_professional_cluster" "professional_cluster" {
  auto_renew_flag = 1
  cluster_name    = "single_zone_cluster"
  product_name    = "PULSAR.P1.MINI2"
  storage_size    = 600
  tags            = {
    "createby" = "terrafrom"
  }
  zone_ids = [
    100004,
  ]

  vpc {
    subnet_id = "subnet-xxxx"
    vpc_id    = "vpc-xxxx"
  }
}
```

Multi-zone Professional Cluster

```hcl
resource "tencentcloud_tdmq_professional_cluster" "professional_cluster" {
  auto_renew_flag = 1
  cluster_name    = "multi_zone_cluster"
  product_name    = "PULSAR.P1.MINI2"
  storage_size    = 200
  tags            = {
    "key"  = "value1"
    "key2" = "value2"
  }
  zone_ids = [
    330001,
    330002,
    330003,
  ]

  vpc {
    subnet_id = "subnet-xxxx"
    vpc_id    = "vpc-xxxx"
  }
}
```

Import

tdmq professional_cluster can be imported using the id, e.g.

```
terraform import tencentcloud_tdmq_professional_cluster.professional_cluster professional_cluster_id
```