Provides a TencentCloud Container Cluster Instance resource.

~> **NOTE:** It has been deprecated and replaced by tencentcloud_kubernetes_scale_worker.

Example Usage

```hcl
resource "tencentcloud_container_cluster_instance" "bar_instance" {
  cpu               = 1
  mem               = 1
  bandwidth         = 1
  bandwidth_type    = "PayByHour"
  require_wan_ip    = 1
  is_vpc_gateway    = 0
  storage_size      = 10
  root_size         = 50
  password          = "Admin12345678"
  cvm_type          = "PayByMonth"
  period            = 1
  zone_id           = 100004
  instance_type     = "CVM.S2"
  mount_target      = "/data"
  docker_graph_path = ""
  subnet_id         = "subnet-abcdedf"
  cluster_id        = "cls-abcdef"
}
```