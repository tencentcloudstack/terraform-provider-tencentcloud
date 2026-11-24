Provides a resource to create an Oceanus (Stream Compute) cluster.

Cluster creation is asynchronous: the provider persists the returned ClusterId
immediately and then polls DescribeClusters until the cluster reports status 2
(running). Only two in-place updates are supported: CU resizing via
ScaleOceanusCluster and prepaid renewal via RenewOceanusCluster. Changing any
other attribute forces the cluster to be recreated.

Note that `region_id`, `zone_id`, `period`, `instance_charge_type`,
`cluster_type`, `renew_flag`, `flink_web_ui_network_access_type`, `cu_memory`
and the VPC lists have no reliable read-back path in DescribeClusters, so drift
on those attributes cannot be detected and their values are preserved from
state. `login_password` is write-only and is never stored in state.

Example Usage

```hcl
resource "tencentcloud_oceanus_oceanus_cluster" "example" {
  cluster_name       = "tf-example-cluster"
  region_id          = 1
  zone_id            = 100001
  login_password     = "Password@123"
  default_cos_bucket = "tf-example-bucket-1250000000"
  cu                 = 12
  remark             = "managed by terraform"

  vpc_descriptions {
    vpc_id    = "vpc-tf1234"
    subnet_id = "subnet-tf1234"
  }

  slave_vpc_descriptions {
    vpc_id    = "vpc-tf1234"
    subnet_id = "subnet-tf1234"
  }
}
```

### Resize an existing cluster

```hcl
resource "tencentcloud_oceanus_oceanus_cluster" "example" {
  cluster_name       = "tf-example-cluster"
  region_id          = 1
  zone_id            = 100001
  login_password     = "Password@123"
  default_cos_bucket = "tf-example-bucket-1250000000"

  # Changing cu triggers ScaleOceanusCluster in place.
  cu = 19

  vpc_descriptions {
    vpc_id    = "vpc-tf1234"
    subnet_id = "subnet-tf1234"
  }
}
```

### Renew a prepaid cluster

```hcl
resource "tencentcloud_oceanus_oceanus_cluster" "example" {
  cluster_name         = "tf-example-cluster"
  region_id            = 1
  zone_id              = 100001
  login_password       = "Password@123"
  default_cos_bucket   = "tf-example-bucket-1250000000"
  cu                   = 12
  instance_charge_type = "PREPAID"
  renew_flag           = "NOTIFY_AND_MANUAL_RENEW"

  # Changing period triggers RenewOceanusCluster in place.
  period = 12

  vpc_descriptions {
    vpc_id    = "vpc-tf1234"
    subnet_id = "subnet-tf1234"
  }
}
```

Import

Oceanus cluster can be imported using the cluster ID, e.g.

```
$ terraform import tencentcloud_oceanus_oceanus_cluster.example cluster-tf123456
```
