---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_instance"
sidebar_current: "docs-tencentcloud-resource-clb_instance"
description: |-
  Provides a resource to create a CLB instance.
---

# tencentcloud_clb_instance

Provides a resource to create a CLB instance.

## Example Usage

```hcl
resource "tencentcloud_clb_instance" "foo" {
  network_type              = "OPEN"
  clb_name                  = "myclb"
  project_id                = 0
  vpc_id                    = "vpc-abcd1234"
  security_groups           = ["sg-o0ek7r93"]
  target_region_info_region = "ap-guangzhou"
  target_region_info_vpc_id = "vpc-abcd1234"

  tags = {
    test = "tf"
  }
}
```

## Argument Reference

The following arguments are supported:

* `clb_name` - (Required) Name of the CLB. The name can only contain Chinese characters, English letters, numbers, underscore and hyphen '-'.
* `network_type` - (Required, ForceNew) Type of CLB instance, and available values include 'OPEN' and 'INTERNAL'.
* `clb_vips` - (Optional) The virtual service address table of the CLB.
* `project_id` - (Optional, ForceNew) ID of the project within the CLB instance, '0' - Default Project.
* `security_groups` - (Optional) Security groups of the CLB instance. Only supports 'OPEN' CLBs.
* `subnet_id` - (Optional, ForceNew) Subnet ID of the CLB. Effective only for CLB within the VPC. Only supports 'INTERNAL' CLBs.
* `tags` - (Optional, ForceNew) The available tags within this CLB.
* `target_region_info_region` - (Optional) Region information of backend services are attached the CLB instance. Only supports 'OPEN' CLBs.
* `target_region_info_vpc_id` - (Optional) Vpc information of backend services are attached the CLB instance. Only supports 'OPEN' CLBs.
* `vpc_id` - (Optional, ForceNew) VPC ID of the CLB.


## Import

CLB instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_instance.foo lb-7a0t6zqb
```

