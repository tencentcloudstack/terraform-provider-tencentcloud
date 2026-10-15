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

* Create: [CreateDedicatedClusterDBInstance](https://cloud.tencent.com/document/api/237/61370) - 创建Mariadb独享集群实例
* Create: [DescribeDBInstances](https://cloud.tencent.com/document/api/237/16184) - 查询实例列表
* Create: [InitDBInstances](https://cloud.tencent.com/document/api/237/16185) - 初始化实例
* Create: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - 批量修改资源关联的标签
* Read: [DescribeDBInstances](https://cloud.tencent.com/document/api/237/16184) - 查询实例列表
* Read: [DescribeResourceTagsByResourceIds](https://cloud.tencent.com/document/api/651/35319) - 批量查看资源关联的标签
* Update: [DescribeFlow](https://cloud.tencent.com/document/api/237/16177) - 查询流程状态
* Update: [ModifyDBInstanceName](https://cloud.tencent.com/document/api/237/16190) - 修改实例名称
* Update: [ModifyDBInstancesProject](https://cloud.tencent.com/document/api/237/16176) - 修改实例所属项目
* Update: [ModifyInstanceNetwork](https://cloud.tencent.com/document/api/237/81949) - 修改实例所属网络
* Update: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - 批量修改资源关联的标签
* Delete: [DescribeDBInstances](https://cloud.tencent.com/document/api/237/16184) - 查询实例列表
* Delete: [IsolateDedicatedDBInstance](https://cloud.tencent.com/document/api/237/82372) - 隔离独享云数据库实例
* Delete: [TerminateDedicatedDBInstance](https://cloud.tencent.com/document/api/237/82555) - 销毁独享云数据库实例

## Import

mariadb dedicatedcluster_db_instance can be imported using the id, e.g.
```
$ terraform import tencentcloud_mariadb_dedicatedcluster_db_instance.dedicatedcluster_db_instance tdsql-050g3fmv
```

