---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_security_group"
sidebar_current: "docs-tencentcloud-resource-cynosdb_security_group"
description: |-
  Provides a resource to create a cynosdb security_group
---

# tencentcloud_cynosdb_security_group

Provides a resource to create a cynosdb security_group

## Example Usage

```hcl
resource "tencentcloud_cynosdb_security_group" "test" {
  cluster_id          = "cynosdbmysql-bws8h88b"
  security_group_ids  = ["sg-baxfiao5"]
  instance_group_type = "RO"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) Cluster id.
* `instance_group_type` - (Required, String, ForceNew) Instance group type. Available values: 
-`HA` - HA group; 
-`RO` - Read-only group;
-`ALL` - HA and RO group.
* `security_group_ids` - (Required, Set: [`String`]) A list of security group IDs to be modified, an array of one or more security group IDs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DescribeClusterDetail](https://cloud.tencent.com/document/api/1003/48086) - This API is used to display cluster details.
* Create: [DescribeClusterInstanceGrps](https://cloud.tencent.com/document/api/1003/48303) - This API is used to query instance groups.
* Create: [DescribeClusters](https://cloud.tencent.com/document/api/1003/48085) - This API is used to describe clusters.
* Create: [ModifyDBInstanceSecurityGroups](https://cloud.tencent.com/document/api/1003/48096) - This API is used to modify the security group bound to the instance.
* Read: [DescribeClusterInstanceGrps](https://cloud.tencent.com/document/api/1003/48303) - This API is used to query instance groups.
* Read: [DescribeDBSecurityGroups](https://cloud.tencent.com/document/api/1003/48305) - This API is used to query instance security group information.
* Update: [DescribeClusterDetail](https://cloud.tencent.com/document/api/1003/48086) - This API is used to display cluster details.
* Update: [DescribeClusterInstanceGrps](https://cloud.tencent.com/document/api/1003/48303) - This API is used to query instance groups.
* Update: [DescribeClusters](https://cloud.tencent.com/document/api/1003/48085) - This API is used to describe clusters.
* Update: [ModifyDBInstanceSecurityGroups](https://cloud.tencent.com/document/api/1003/48096) - This API is used to modify the security group bound to the instance.
* Delete: [DescribeClusterDetail](https://cloud.tencent.com/document/api/1003/48086) - This API is used to display cluster details.
* Delete: [DescribeClusterInstanceGrps](https://cloud.tencent.com/document/api/1003/48303) - This API is used to query instance groups.
* Delete: [DescribeClusters](https://cloud.tencent.com/document/api/1003/48085) - This API is used to describe clusters.
* Delete: [DescribeDBSecurityGroups](https://cloud.tencent.com/document/api/1003/48305) - This API is used to query instance security group information.
* Delete: [DisassociateSecurityGroups](https://cloud.tencent.com/document/api/1003/70117) - This API is used to unbind cloud resources from security groups in batches.

## Import

cynosdb security_group can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_security_group.security_group ${cluster_id}#${instance_group_type}
```

