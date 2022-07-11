---
subcategory: "Tencent Container Registry(TCR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_vpc_attachment"
sidebar_current: "docs-tencentcloud-resource-tcr_vpc_attachment"
description: |-
  Use this resource to create tcr vpc attachment to manage access of internal endpoint.
---

# tencentcloud_tcr_vpc_attachment

Use this resource to create tcr vpc attachment to manage access of internal endpoint.

## Example Usage

```hcl
resource "tencentcloud_tcr_vpc_attachment" "foo" {
  instance_id = "cls-satg5125"
  vpc_id      = "vpc-asg3sfa3"
  subnet_id   = "subnet-1uwh63so"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) ID of the TCR instance.
* `subnet_id` - (Required, String, ForceNew) ID of subnet.
* `vpc_id` - (Required, String, ForceNew) ID of VPC.
* `enable_public_domain_dns` - (Optional, Bool) Whether to enable public domain dns. Default value is `false`.
* `enable_vpc_domain_dns` - (Optional, Bool) Whether to enable vpc domain dns. Default value is `false`.
* `region_id` - (Optional, Int) ID of region. Conflict with region_name, can not be set at the same time.
* `region_name` - (Optional, String) Name of region. Conflict with region_id, can not be set at the same time.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `access_ip` - IP address of the internal access.
* `status` - Status of the internal access.


## Import

tcr vpc attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_tcr_vpc_attachment.foo tcrId#vpcId#subnetId
```

