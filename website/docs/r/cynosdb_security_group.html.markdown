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

* Create: [DescribeClusterDetail](https://www.tencentcloud.com/search/?wd=DescribeClusterDetail&type=document&lang=en)
* Create: [DescribeClusterInstanceGrps](https://www.tencentcloud.com/search/?wd=DescribeClusterInstanceGrps&type=document&lang=en)
* Create: [DescribeClusters](https://www.tencentcloud.com/search/?wd=DescribeClusters&type=document&lang=en)
* Create: [ModifyDBInstanceSecurityGroups](https://www.tencentcloud.com/search/?wd=ModifyDBInstanceSecurityGroups&type=document&lang=en)
* Read: [DescribeClusterInstanceGrps](https://www.tencentcloud.com/search/?wd=DescribeClusterInstanceGrps&type=document&lang=en)
* Read: [DescribeDBSecurityGroups](https://www.tencentcloud.com/search/?wd=DescribeDBSecurityGroups&type=document&lang=en)
* Update: [DescribeClusterDetail](https://www.tencentcloud.com/search/?wd=DescribeClusterDetail&type=document&lang=en)
* Update: [DescribeClusterInstanceGrps](https://www.tencentcloud.com/search/?wd=DescribeClusterInstanceGrps&type=document&lang=en)
* Update: [DescribeClusters](https://www.tencentcloud.com/search/?wd=DescribeClusters&type=document&lang=en)
* Update: [ModifyDBInstanceSecurityGroups](https://www.tencentcloud.com/search/?wd=ModifyDBInstanceSecurityGroups&type=document&lang=en)
* Delete: [DescribeClusterDetail](https://www.tencentcloud.com/search/?wd=DescribeClusterDetail&type=document&lang=en)
* Delete: [DescribeClusterInstanceGrps](https://www.tencentcloud.com/search/?wd=DescribeClusterInstanceGrps&type=document&lang=en)
* Delete: [DescribeClusters](https://www.tencentcloud.com/search/?wd=DescribeClusters&type=document&lang=en)
* Delete: [DescribeDBSecurityGroups](https://www.tencentcloud.com/search/?wd=DescribeDBSecurityGroups&type=document&lang=en)
* Delete: [DisassociateSecurityGroups](https://www.tencentcloud.com/search/?wd=DisassociateSecurityGroups&type=document&lang=en)

## Import

cynosdb security_group can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_security_group.security_group ${cluster_id}#${instance_group_type}
```

