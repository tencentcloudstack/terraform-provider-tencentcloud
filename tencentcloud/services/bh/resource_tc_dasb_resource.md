Provides a resource to create a dasb resource

Example Usage

Create a standard version instance

```hcl
resource "tencentcloud_dasb_resource" "example" {
  deploy_region     = "ap-guangzhou"
  deploy_zone       = "ap-guangzhou-6"
  vpc_id            = "vpc-fmz6l9nz"
  subnet_id         = "subnet-g7jhwhi2"
  vpc_cidr_block    = "10.35.0.0/16"
  cidr_block        = "10.35.20.0/24"
  resource_edition  = "standard"
  resource_node     = 50
  time_unit         = "m"
  time_span         = 1
  auto_renew_flag   = 1
  package_bandwidth = 1
}
```

Create a professional instance

```hcl
resource "tencentcloud_dasb_resource" "example" {
  deploy_region     = "ap-guangzhou"
  deploy_zone       = "ap-guangzhou-6"
  vpc_id            = "vpc-fmz6l9nz"
  subnet_id         = "subnet-g7jhwhi2"
  vpc_cidr_block    = "10.35.0.0/16"
  cidr_block        = "10.35.20.0/24"
  resource_edition  = "pro"
  resource_node     = 50
  time_unit         = "m"
  time_span         = 1
  auto_renew_flag   = 1
  package_bandwidth = 1
}
```

Import

dasb resource can be imported using the id, e.g.

```
terraform import tencentcloud_dasb_resource.example bh-saas-kgckynrt
```
