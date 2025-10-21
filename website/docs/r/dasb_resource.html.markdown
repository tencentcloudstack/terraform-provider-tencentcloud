---
subcategory: "Bastion Host(BH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dasb_resource"
sidebar_current: "docs-tencentcloud-resource-dasb_resource"
description: |-
  Provides a resource to create a dasb resource
---

# tencentcloud_dasb_resource

Provides a resource to create a dasb resource

## Example Usage

### Create a standard version instance

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

### Create a professional instance

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

## Argument Reference

The following arguments are supported:

* `auto_renew_flag` - (Required, Int) Automatic renewal. 1 is auto renew flag, 0 is not.
* `cidr_block` - (Required, String) Subnet segments that require service activation.
* `deploy_region` - (Required, String) Deploy region.
* `deploy_zone` - (Required, String) Deploy zone.
* `resource_edition` - (Required, String) Resource type.Value:standard/pro.
* `resource_node` - (Required, Int) Number of resource nodes.
* `subnet_id` - (Required, String) Deploy resource subnetId.
* `vpc_cidr_block` - (Required, String) The network segment corresponding to the VPC that requires service activation.
* `vpc_id` - (Required, String) Deploy resource vpcId.
* `package_bandwidth` - (Optional, Int) Number of bandwidth expansion packets (4M), The set value is an integer multiple of 4.
* `time_span` - (Optional, Int) Billing time. This field is mandatory, with a minimum value of 1.
* `time_unit` - (Optional, String) Billing cycle, only support m: month. This field is mandatory, fill in m.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

dasb resource can be imported using the id, e.g.

```
terraform import tencentcloud_dasb_resource.example bh-saas-kgckynrt
```

