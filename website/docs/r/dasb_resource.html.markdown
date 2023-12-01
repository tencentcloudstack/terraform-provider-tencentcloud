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

## Argument Reference

The following arguments are supported:

* `auto_renew_flag` - (Required, Int) Automatic renewal. 1 is auto renew flag, 0 is not.
* `deploy_region` - (Required, String) Deploy region.
* `resource_edition` - (Required, String) Resource type.Value:standard/pro.
* `resource_node` - (Required, Int) Number of resource nodes.
* `subnet_id` - (Required, String) Deploy resource subnetId.
* `time_span` - (Required, Int) Billing time.
* `time_unit` - (Required, String) Billing cycle, only support m: month.
* `vpc_id` - (Required, String) Deploy resource vpcId.
* `deploy_zone` - (Optional, String) Deploy zone.
* `package_bandwidth` - (Optional, Int) Number of bandwidth expansion packets (4M).
* `package_node` - (Optional, Int) Number of authorized point extension packages (50 points).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

dasb resource can be imported using the id, e.g.

```
terraform import tencentcloud_dasb_resource.example bh-saas-kk5rabk0
```

