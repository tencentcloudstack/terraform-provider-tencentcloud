---
subcategory: "TencentDB for MariaDB(MariaDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mariadb_dedicatedcluster_db_instance"
sidebar_current: "docs-tencentcloud-resource-mariadb_dedicatedcluster_db_instance"
description: |-
  Provides a resource to create a mariadb dedicatedcluster_db_instance
---

# tencentcloud_mariadb_dedicatedcluster_db_instance

Provides a resource to create a mariadb dedicatedcluster_db_instance

## Example Usage

```hcl
resource "tencentcloud_mariadb_dedicatedcluster_db_instance" "dedicatedcluster_db_instance" {
  goods_num     = 1
  memory        = 2
  storage       = 10
  cluster_id    = "dbdc-24odnuhr"
  vpc_id        = "vpc-ii1jfbhl"
  subnet_id     = "subnet-3ku415by"
  db_version_id = "8.0"
  instance_name = "cluster-mariadb-test-1"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) dedicated cluster id.
* `goods_num` - (Required, Int) number of instance.
* `memory` - (Required, Int) instance memory.
* `storage` - (Required, Int) instance disk storage.
* `db_version_id` - (Optional, String) db engine version, default to 0.
* `instance_name` - (Optional, String) name of this instance.
* `project_id` - (Optional, Int) project id.
* `subnet_id` - (Optional, String) subnet id, it&amp;#39;s required when vpcId is set.
* `tags` - (Optional, Map) Tag description list.
* `vip` - (Optional, String) vip.
* `vpc_id` - (Optional, String) vpc id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateDedicatedClusterDBInstance](https://cloud.tencent.com/document/api/237/61370)
* Create: [DescribeDBInstances](https://www.tencentcloud.com/document/product/237/16184) - This API is used to query the TencentDB instance list. It supports filtering instances by project ID, instance ID, private address, and instance name.
* Create: [InitDBInstances](https://www.tencentcloud.com/document/product/237/16185)
* Create: [ModifyResourceTags](https://www.tencentcloud.com/document/product/651/33012) - This API is used to modify all tags associated with a resource.
* Read: [DescribeDBInstances](https://www.tencentcloud.com/document/product/237/16184) - This API is used to query the TencentDB instance list. It supports filtering instances by project ID, instance ID, private address, and instance name.
* Read: [DescribeResourceTagsByResourceIds](https://www.tencentcloud.com/document/product/651/33022) - This API is used to query the tag key-value pairs associated with an existing resource.
* Update: [DescribeFlow](https://www.tencentcloud.com/document/product/237/16177)
* Update: [ModifyDBInstanceName](https://www.tencentcloud.com/document/product/237/16190)
* Update: [ModifyDBInstancesProject](https://www.tencentcloud.com/document/product/237/16176) - This API is used to modify the project to which TencentDB instances belong.
* Update: [ModifyInstanceNetwork](https://www.tencentcloud.com/document/product/237/51318) - This API is used to modify instance network.
* Update: [ModifyResourceTags](https://www.tencentcloud.com/document/product/651/33012) - This API is used to modify all tags associated with a resource.
* Delete: [DescribeDBInstances](https://www.tencentcloud.com/document/product/237/16184) - This API is used to query the TencentDB instance list. It supports filtering instances by project ID, instance ID, private address, and instance name.
* Delete: [IsolateDedicatedDBInstance](https://www.tencentcloud.com/document/product/237/51453) - This API is used to isolate a dedicated TencentDB instance.
* Delete: [TerminateDedicatedDBInstance](https://www.tencentcloud.com/document/product/237/51452) - This API is used to terminate the isolated dedicated TencentDB instance.

## Import

mariadb dedicatedcluster_db_instance can be imported using the id, e.g.
```
$ terraform import tencentcloud_mariadb_dedicatedcluster_db_instance.dedicatedcluster_db_instance tdsql-050g3fmv
```

