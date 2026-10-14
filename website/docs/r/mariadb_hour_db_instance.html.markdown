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

* Create: [CreateHourDBInstance](https://www.tencentcloud.com/document/product/237/45098) - This API is used to create a pay-as-you-go TencentDB for MariaDB instance.
* Create: [DescribeDBInstances](https://www.tencentcloud.com/document/product/237/16184) - This API is used to query the TencentDB instance list. It supports filtering instances by project ID, instance ID, private address, and instance name.
* Create: [InitDBInstances](https://www.tencentcloud.com/document/product/237/16185)
* Create: [ModifyResourceTags](https://www.tencentcloud.com/document/product/651/33012) - This API is used to modify all tags associated with a resource.
* Read: [DescribeDBInstanceDetail](https://www.tencentcloud.com/document/product/237/54563) - This API is used to query the details of a specified instance.
* Read: [DescribeDBParameters](https://www.tencentcloud.com/document/product/237/16154) - This API is used to get the current parameter settings of a database.
* Read: [DescribeResourceTagsByResourceIds](https://www.tencentcloud.com/document/product/651/33022) - This API is used to query the tag key-value pairs associated with an existing resource.
* Update: [DescribeFlow](https://www.tencentcloud.com/document/product/237/16177)
* Update: [ModifyDBInstanceName](https://www.tencentcloud.com/document/product/237/16190)
* Update: [ModifyDBInstancesProject](https://www.tencentcloud.com/document/product/237/16176) - This API is used to modify the project to which TencentDB instances belong.
* Update: [ModifyDBParameters](https://www.tencentcloud.com/document/product/237/16153) - This API is used to modify database parameters.
* Update: [ModifyInstanceNetwork](https://www.tencentcloud.com/document/product/237/51318) - This API is used to modify instance network.
* Update: [ModifyResourceTags](https://www.tencentcloud.com/document/product/651/33012) - This API is used to modify all tags associated with a resource.
* Delete: [DestroyHourDBInstance](https://www.tencentcloud.com/document/product/237/40958) - This API is used to terminate a pay-as-you-go TencentDB for MariaDB instance.

## Import

MariaDB hour db instance can be imported using the id, e.g.
```
terraform import tencentcloud_mariadb_hour_db_instance.example tdsql-kjqih9nn
```

