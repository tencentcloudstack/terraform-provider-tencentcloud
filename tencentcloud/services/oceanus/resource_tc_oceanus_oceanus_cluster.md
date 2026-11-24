Provides a resource to create an Oceanus (Stream Compute) cluster.

Cluster creation is asynchronous: the provider persists the returned ClusterId
immediately and then polls DescribeClusters until the cluster reports status 2
(running). In-place updates are supported only for CU resizing (via
ScaleOceanusCluster) and prepaid renewal (via RenewOceanusCluster); changing
any other attribute forces the cluster to be recreated. Note that
`region_id`, `zone_id`, `login_password`, `period`,
`instance_charge_type`, `cluster_type`, `renew_flag`,
`flink_web_ui_network_access_type`, `cu_memory` and the VPC lists have no
reliable read-back path in DescribeClusters, so drift on those attributes
cannot be detected and their values are preserved from state.

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

Import

Oceanus cluster can be imported using the cluster ID, e.g.

```
$ terraform import tencentcloud_oceanus_oceanus_cluster.example cluster-tf123456
```
