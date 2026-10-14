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

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DescribeInstances](https://cloud.tencent.com/document/api/1141/41569) - This API is used to query the instance information.
* Create: [DescribeInternalEndpoints](https://cloud.tencent.com/document/api/1141/50734) - This API is used to query the VPC URLs for private network access to an instance.
* Create: [ManageInternalEndpoint](https://cloud.tencent.com/document/api/1141/50733) - This API is used to manage VPC URLs for private network access to an instance.
* Read: [DescribeInstances](https://cloud.tencent.com/document/api/1141/41569) - This API is used to query the instance information.
* Read: [DescribeInternalEndpoints](https://cloud.tencent.com/document/api/1141/50734) - This API is used to query the VPC URLs for private network access to an instance.
* Delete: [DescribeInstances](https://cloud.tencent.com/document/api/1141/41569) - This API is used to query the instance information.
* Delete: [DescribeInternalEndpoints](https://cloud.tencent.com/document/api/1141/50734) - This API is used to query the VPC URLs for private network access to an instance.
* Delete: [ManageInternalEndpoint](https://cloud.tencent.com/document/api/1141/50733) - This API is used to manage VPC URLs for private network access to an instance.

## Import

tcr vpc attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_tcr_vpc_attachment.foo instance_id#vpc_id#subnet_id
```

