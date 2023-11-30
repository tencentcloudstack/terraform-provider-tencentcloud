Provides a TencentCloud Container Cluster resource.

~> **NOTE:** It has been deprecated and replaced by tencentcloud_kubernetes_cluster.

Example Usage

```hcl
resource "tencentcloud_container_cluster" "foo" {
  cluster_name                 = "terraform-acc-test"
  cpu                          = 1
  mem                          = 1
  os_name                      = "ubuntu16.04.1 LTSx86_64"
  bandwidth                    = 1
  bandwidth_type               = "PayByHour"
  require_wan_ip               = 1
  subnet_id                    = "subnet-abcdabc"
  is_vpc_gateway               = 0
  storage_size                 = 0
  root_size                    = 50
  goods_num                    = 1
  password                     = "Admin12345678"
  vpc_id                       = "vpc-abcdabc"
  cluster_cidr                 = "10.0.2.0/24"
  ignore_cluster_cidr_conflict = 0
  cvm_type                     = "PayByHour"
  cluster_desc                 = "foofoofoo"
  period                       = 1
  zone_id                      = 100004
  instance_type                = "S2.SMALL1"
  mount_target                 = ""
  docker_graph_path            = ""
  instance_name                = "bar-vm"
  cluster_version              = "1.7.8"
}
```