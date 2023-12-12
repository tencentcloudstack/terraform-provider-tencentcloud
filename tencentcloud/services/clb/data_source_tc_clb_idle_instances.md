Use this data source to query detailed information of clb idle_loadbalancers

Example Usage

```hcl
data "tencentcloud_clb_idle_instances" "idle_instance" {
  load_balancer_region = "ap-guangzhou"
}
```