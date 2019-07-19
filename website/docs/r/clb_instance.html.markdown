---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_instance"
sidebar_current: "docs-tencentcloud-resource-clb_instance"
description: |-
  Provide a resource to create a CLB instance.
---

# tencentcloud_clb_instance

Provide a resource to create a CLB instance.

## Example Usage

```hcl
resource "tencentcloud_clb_listener" "clblab" {
  network_type              = "OPEN"
  clb_name                  = "myclb"
  project_id                = "Default Project"
  vpc_id                    = "vpc-abcd1234"
  subnet_id                 = "subnet-0agspqdn"
  tags                      = "mytags"
  security_groups           = ["sg-o0ek7r93"]
  target_region_info_region = "ap-guangzhou"
  target_region_info_vpc_id = "vpc-abcd1234"
}
```

## Argument Reference

The following arguments are supported:

* `clb_name` - (Required) Name of the CLB to be queried. The name can only contain Chinese characters, English letters, numbers, underscore and hyphen '-'
* `network_type` - (Required, ForceNew) Type of CLB instance, and available values include 'OPEN' and 'INTERNAL'.
* `project_id` - (Optional, ForceNew) ID of the project to which the instance belongs.
* `security_groups` - (Optional) Security groups to which a CLB instance belongs.
* `subnet_id` - (Optional, ForceNew) ID of the subnet within this VPC. The VIP of the intranet CLB instance will be generated from this subnet
* `target_region_info_region` - (Optional) Region information of backend service are attached the CLB instance.
* `target_region_info_vpc` - (Optional) Vpc Id information of backend service are attached the CLB instance.
* `vpc_id` - (Optional, ForceNew) ID of the subnet within this VPC. The VIP of the intranet CLB instance will be generated from this subnet


## Import

CLB instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb.instance lb-7a0t6zqb
```

