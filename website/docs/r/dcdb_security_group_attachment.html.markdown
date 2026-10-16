---
subcategory: "TDSQL for MySQL(DCDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dcdb_security_group_attachment"
sidebar_current: "docs-tencentcloud-resource-dcdb_security_group_attachment"
description: |-
  Provides a resource to create a dcdb security_group_attachment
---

# tencentcloud_dcdb_security_group_attachment

Provides a resource to create a dcdb security_group_attachment

## Example Usage

```hcl
resource "tencentcloud_dcdb_security_group_attachment" "security_group_attachment" {
  security_group_id = ""
  instance_id       = ""
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) attached instance id.
* `security_group_id` - (Required, String, ForceNew) security group id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [AssociateSecurityGroups](https://cloud.tencent.com/document/api/557/51782)
* Read: [DescribeDBSecurityGroups](https://cloud.tencent.com/document/api/557/51781)
* Delete: [DisassociateSecurityGroups](https://cloud.tencent.com/document/api/557/51779)

## Import

dcdb security_group_attachment can be imported using the id, e.g.
```
$ terraform import tencentcloud_dcdb_security_group_attachment.security_group_attachment securityGroupAttachment_id
```

