---
subcategory: "Tencent Container Registry(TCR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_vpc_attachment"
sidebar_current: "docs-tencentcloud-resource-tcr_vpc_attachment"
description: |-
  Use this resource to attach tcr instance with the vpc and subnet network.
---

# tencentcloud_tcr_vpc_attachment

Use this resource to attach tcr instance with the vpc and subnet network.

## Example Usage

### Attach a tcr instance with vpc resource

```hcl
locals {
  vpc_id    = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  subnet_id = data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id
  tcr_id    = tencentcloud_tcr_instance.example.id
}

data "tencentcloud_vpc_subnets" "vpc" {
  is_default        = true
  availability_zone = var.availability_zone
}

data "tencentcloud_security_groups" "sg" {
  name = "default"
}

resource "tencentcloud_tcr_instance" "example" {
  name          = "tf-example-tcr"
  instance_type = "basic"
  delete_bucket = true
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tcr_vpc_attachment" "foo" {
  instance_id = local.tcr_id
  vpc_id      = local.vpc_id
  subnet_id   = local.subnet_id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) ID of the TCR instance.
* `subnet_id` - (Required, String, ForceNew) ID of subnet.
* `vpc_id` - (Required, String, ForceNew) ID of VPC.
* `enable_public_domain_dns` - (Optional, Bool) Whether to enable public domain dns. Default value is `false`.
* `enable_vpc_domain_dns` - (Optional, Bool) Whether to enable vpc domain dns. Default value is `false`.
* `region_id` - (Optional, Int, **Deprecated**) this argument was deprecated, use `region_name` instead. ID of region. Conflict with region_name, can not be set at the same time.
* `region_name` - (Optional, String) Name of region. Conflict with region_id, can not be set at the same time.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `access_ip` - IP address of the internal access.
* `status` - Status of the internal access.


## Import

tcr vpc attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_tcr_vpc_attachment.foo instance_id#vpc_id#subnet_id
```

