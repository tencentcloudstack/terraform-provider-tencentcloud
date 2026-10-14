---
subcategory: "TencentDB for MariaDB(MariaDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mariadb_security_groups"
sidebar_current: "docs-tencentcloud-resource-mariadb_security_groups"
description: |-
  Provides a resource to create a mariadb security_groups
---

# tencentcloud_mariadb_security_groups

Provides a resource to create a mariadb security_groups

~> **NOTE:** If you use this resource, please do not use security_group_ids in tencentcloud_mariadb_instance resource

## Example Usage

```hcl
resource "tencentcloud_mariadb_security_groups" "security_groups" {
  instance_id       = "tdsql-4pzs5b67"
  security_group_id = "sg-7kpsbxdb"
  product           = "mariadb"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) instance id.
* `product` - (Required, String) product name, fixed to mariadb.
* `security_group_id` - (Required, String) security group id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [AssociateSecurityGroups](https://www.tencentcloud.com/document/product/237/39450) - This API is used to associate security groups with Tencent Cloud resources in batches.
* Read: [DescribeDBSecurityGroups](https://www.tencentcloud.com/document/product/237/39449) - This API is used to query the security group details of an instance.
* Update: [ModifyDBInstanceSecurityGroups](https://www.tencentcloud.com/document/product/237/39446)
* Delete: [DisassociateSecurityGroups](https://www.tencentcloud.com/document/product/237/39447) - This API is used to unassociate security groups from instances in batches.

## Import

mariadb security_groups can be imported using the id, e.g.
```
$ terraform import tencentcloud_mariadb_security_groups.security_groups tdsql-4pzs5b67#sg-7kpsbxdb#mariadb
```

