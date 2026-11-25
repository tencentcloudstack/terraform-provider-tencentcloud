Provides a resource to create an Oceanus (Stream Compute) cluster.

~> **NOTE:** `login_password` is write-only and is never stored in state.

Example Usage

```hcl
resource "tencentcloud_oceanus_cluster" "example" {
  cluster_name                     = "tf-test4"
  cluster_type                     = ""
  cu                               = 12
  cu_memory                        = 4
  default_cos_bucket               = "keep-tf-test-1308726196"
  flink_web_ui_network_access_type = "NetworkAccess_EXTERNAL"
  instance_charge_type             = "PREPAID"
  login_password                   = "Password@123"
  period                           = 1
  region                           = "ap-guangzhou"
  zone                             = "ap-guangzhou-7"
  remark                           = "测试"
  renew_flag                       = "NOTIFY_AND_MANUAL_RENEW"
  vpc_descriptions = [{
    subnet_id = "subnet-4sp1mbni"
    vpc_id    = "vpc-2trcrynd"
  }]
}
```

Import

Oceanus cluster can be imported using the cluster ID, e.g.

```
terraform import tencentcloud_oceanus_cluster.example cluster-tf123456
```
