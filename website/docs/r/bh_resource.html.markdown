---
subcategory: "Bastion Host(BH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_bh_resource"
sidebar_current: "docs-tencentcloud-resource-bh_resource"
description: |-
  Provides a resource to create a BH resource
---

# tencentcloud_bh_resource

Provides a resource to create a BH resource

~> **NOTE:** Currently, executing the `terraform destroy` command to delete this resource is not supported. If you need to destroy it, please contact Tencent Cloud BH through a ticket.

## Example Usage

```hcl
resource "tencentcloud_bh_resource" "example" {
  deploy_region    = "ap-guangzhou"
  vpc_id           = "vpc-q1of50wz"
  subnet_id        = "subnet-7uhvm46o"
  resource_edition = "standard"
  resource_node    = 20
  time_unit        = "m"
  time_span        = "1"
  pay_mode         = 1
  auto_renew_flag  = 1
  deploy_zone      = "ap-guangzhou-6"
  cidr_block       = "192.168.11.0/24"
  vpc_cidr_block   = "192.168.0.0/16"
  intranet_access  = 1
  external_access  = 1
}
```

## Argument Reference

The following arguments are supported:

* `auto_renew_flag` - (Required, Int) Auto-renewal.
* `cidr_block` - (Required, String, ForceNew) CIDR block of the bastion host.
* `deploy_region` - (Required, String, ForceNew) Deployment region.
* `deploy_zone` - (Required, String, ForceNew) Deployment zone.
* `pay_mode` - (Required, Int, ForceNew) Billing mode, 1 for prepaid.
* `resource_edition` - (Required, String) Resource type. Values: standard/pro.
* `resource_node` - (Required, Int) Number of resource nodes.
* `subnet_id` - (Required, String, ForceNew) Subnet ID for deploying the bastion host.
* `time_span` - (Required, Int, ForceNew) Billing duration.
* `time_unit` - (Required, String, ForceNew) Billing cycle.
* `vpc_cidr_block` - (Required, String, ForceNew) The network segment corresponding to the VPC that needs to activate the service.
* `vpc_id` - (Required, String, ForceNew) VPC ID for deploying the bastion host.
* `client_access` - (Optional, Int, ForceNew) 0 - Disable client access to the bastion host; 1 - Enable client access to the bastion host.
* `external_access` - (Optional, Int) 0 - Disable public network access to the bastion host; 1 - Enable public network access to the bastion host.
* `intranet_access` - (Optional, Int) 0 - Disable internal network access bastion host; 1 - Enable internal network access bastion host.
* `share_clb` - (Optional, Int, ForceNew) Whether to share CLB, 0: not shared, 1: shared.
* `trial` - (Optional, Int, ForceNew) 0 for non-trial version, 1 for trial version.
* `web_access` - (Optional, Int, ForceNew) 0 - Disable web access bastion host; 1 - Enable web access bastion host.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `resource_id` - Resource instance ID.


