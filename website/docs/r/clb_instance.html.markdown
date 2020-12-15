---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_instance"
sidebar_current: "docs-tencentcloud-resource-clb_instance"
description: |-
  Provides a resource to create a CLB instance.
---

# tencentcloud_clb_instance

Provides a resource to create a CLB instance.

## Example Usage

INTERNAL CLB

```hcl
resource "tencentcloud_clb_instance" "internal_clb" {
  network_type = "INTERNAL"
  clb_name     = "myclb"
  project_id   = 0
  vpc_id       = "vpc-7007ll7q"
  subnet_id    = "subnet-12rastkr"

  tags = {
    test = "tf"
  }
}
```

OPEN CLB

```hcl
resource "tencentcloud_clb_instance" "open_clb" {
  network_type              = "OPEN"
  clb_name                  = "myclb"
  project_id                = 0
  vpc_id                    = "vpc-da7ffa61"
  security_groups           = ["sg-o0ek7r93"]
  target_region_info_region = "ap-guangzhou"
  target_region_info_vpc_id = "vpc-da7ffa61"

  tags = {
    test = "tf"
  }
}
```

## Argument Reference

The following arguments are supported:

* `clb_name` - (Required) Name of the CLB. The name can only contain Chinese characters, English letters, numbers, underscore and hyphen '-'.
* `network_type` - (Required, ForceNew) Type of CLB instance. Valid values: `OPEN` and `INTERNAL`.
* `project_id` - (Optional, ForceNew) ID of the project within the CLB instance, `0` - Default Project.
* `security_groups` - (Optional) Security groups of the CLB instance. Only supports `OPEN` CLBs.
* `subnet_id` - (Optional, ForceNew) Subnet ID of the CLB. Effective only for CLB within the VPC. Only supports `INTERNAL` CLBs.
* `tags` - (Optional) The available tags within this CLB.
* `target_region_info_region` - (Optional) Region information of backend services are attached the CLB instance. Only supports `OPEN` CLBs.
* `target_region_info_vpc_id` - (Optional) Vpc information of backend services are attached the CLB instance. Only supports `OPEN` CLBs.
* `vpc_id` - (Optional, ForceNew) VPC ID of the CLB.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `clb_vips` - The virtual service address table of the CLB.


## Import

CLB instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_instance.foo lb-7a0t6zqb
```

