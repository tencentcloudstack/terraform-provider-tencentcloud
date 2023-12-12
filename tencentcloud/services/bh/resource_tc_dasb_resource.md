Provides a resource to create a dasb resource

Example Usage

```hcl
resource "tencentcloud_dasb_resource" "example" {
  deploy_region     = "ap-guangzhou"
  vpc_id            = "vpc-q1of50wz"
  subnet_id         = "subnet-7uhvm46o"
  resource_edition  = "standard"
  resource_node     = 2
  time_unit         = "m"
  time_span         = 1
  auto_renew_flag   = 1
  deploy_zone       = "ap-guangzhou-6"
  package_bandwidth = 10
  package_node      = 50
}
```

Import

dasb resource can be imported using the id, e.g.

```
terraform import tencentcloud_dasb_resource.example bh-saas-kk5rabk0
```