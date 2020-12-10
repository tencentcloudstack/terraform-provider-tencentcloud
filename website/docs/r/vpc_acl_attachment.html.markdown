---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_acl_attachment"
sidebar_current: "docs-tencentcloud-resource-vpc_acl_attachment"
description: |-
  Provide a resource to attach an existing subnet to Network ACL.
---

# tencentcloud_vpc_acl_attachment

Provide a resource to attach an existing subnet to Network ACL.

## Example Usage

```hcl
data "tencentcloud_vpc_instances" "id_instances" {
}
resource "tencentcloud_vpc_acl" "foo" {
  vpc_id = data.tencentcloud_vpc_instances.id_instances.instance_list.0.vpc_id
  name   = "test_acl"
  ingress = [
    "ACCEPT#192.168.1.0/24#800#TCP",
    "ACCEPT#192.168.1.0/24#800-900#TCP",
  ]
  egress = [
    "ACCEPT#192.168.1.0/24#800#TCP",
    "ACCEPT#192.168.1.0/24#800-900#TCP",
  ]
}

resource "tencentcloud_vpc_acl_attachment" "attachment" {
  acl_id    = tencentcloud_vpc_acl.foo.id
  subnet_id = data.tencentcloud_vpc_instances.id_instances.instance_list[0].subnet_ids[0]
}
```

## Argument Reference

The following arguments are supported:

* `acl_id` - (Required, ForceNew) ID of the attached ACL.
* `subnet_id` - (Required, ForceNew) The Subnet instance ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Acl attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_vpc_acl_attachment.attachment acl-eotx5qsg#subnet-91x0geu6
```

