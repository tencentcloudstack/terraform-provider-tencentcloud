---
subcategory: "TencentDB for MariaDB(MariaDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mariadb_hour_db_instance"
sidebar_current: "docs-tencentcloud-resource-mariadb_hour_db_instance"
description: |-
  Provides a resource to create a MariaDB hour db instance
---

# tencentcloud_mariadb_hour_db_instance

Provides a resource to create a MariaDB hour db instance

## Example Usage

### Create with default init params

```hcl
resource "tencentcloud_mariadb_hour_db_instance" "example" {
  instance_name = "tf-example"
  memory        = 4
  node_count    = 2
  storage       = 100
  vpc_id        = "vpc-i5yyodl9"
  subnet_id     = "subnet-d4umunpy"
  vip           = "10.0.0.8"
  zones = [
    "ap-guangzhou-6",
    "ap-guangzhou-7",
  ]

  tags = {
    createdBy = "Terraform"
  }
}
```

### Create with custom init params

```hcl
resource "tencentcloud_mariadb_hour_db_instance" "example" {
  db_version_id = "5.7"
  instance_name = "tf-example"
  memory        = 2
  node_count    = 2
  storage       = 100
  vpc_id        = "vpc-i5yyodl9"
  subnet_id     = "subnet-d4umunpy"
  vip           = "10.0.0.8"
  zones = [
    "ap-guangzhou-6",
    "ap-guangzhou-7",
  ]

  init_params {
    param = "character_set_server"
    value = "utf8"
  }

  init_params {
    param = "lower_case_table_names"
    value = "1"
  }

  tags = {
    createdBy = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `memory` - (Required, Int) instance memory.
* `node_count` - (Required, Int) number of node for instance.
* `storage` - (Required, Int) instance disk storage.
* `zones` - (Required, List: [`String`]) available zone of instance.
* `db_version_id` - (Optional, String) db engine version, default to 10.1.9.
* `init_params` - (Optional, List) parameter list. This interface's optional values include: `character_set_server` (character set, required), `lower_case_table_names` (table name case sensitivity, required, 0 - sensitive; 1 - insensitive), `innodb_page_size` (innodb data page, default 16K), `sync_mode` (sync mode: 0 - async; 1 - strong sync; 2 - strong sync degradable, default is strong sync degradable).
* `instance_name` - (Optional, String) name of this instance.
* `project_id` - (Optional, Int) project id.
* `subnet_id` - (Optional, String) subnet id, it&amp;#39;s required when vpcId is set.
* `tags` - (Optional, Map) Tag description list.
* `vip` - (Optional, String) vip.
* `vpc_id` - (Optional, String) vpc id.

The `init_params` object supports the following:

* `param` - (Required, String) parameter name.
* `value` - (Required, String) parameter value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateHourDBInstance](https://cloud.tencent.com/document/api/237/67878) - 创建MariaDB按量计费实例
* Create: [DescribeDBInstances](https://cloud.tencent.com/document/api/237/16184) - 查询实例列表
* Create: [InitDBInstances](https://cloud.tencent.com/document/api/237/16185) - 初始化实例
* Create: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - 批量修改资源关联的标签
* Read: [DescribeDBInstanceDetail](https://cloud.tencent.com/document/api/237/89390) - 查询实例详细信息
* Read: [DescribeDBParameters](https://cloud.tencent.com/document/api/237/16154) - 查看数据库参数
* Read: [DescribeResourceTagsByResourceIds](https://cloud.tencent.com/document/api/651/35319) - 批量查看资源关联的标签
* Update: [DescribeFlow](https://cloud.tencent.com/document/api/237/16177) - 查询流程状态
* Update: [ModifyDBInstanceName](https://cloud.tencent.com/document/api/237/16190) - 修改实例名称
* Update: [ModifyDBInstancesProject](https://cloud.tencent.com/document/api/237/16176) - 修改实例所属项目
* Update: [ModifyDBParameters](https://cloud.tencent.com/document/api/237/16153) - 修改数据库参数
* Update: [ModifyInstanceNetwork](https://cloud.tencent.com/document/api/237/81949) - 修改实例所属网络
* Update: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - 批量修改资源关联的标签
* Delete: [DestroyHourDBInstance](https://cloud.tencent.com/document/api/237/56487) - 销毁MariaDB按量计费实例

## Import

MariaDB hour db instance can be imported using the id, e.g.
```
terraform import tencentcloud_mariadb_hour_db_instance.example tdsql-kjqih9nn
```

