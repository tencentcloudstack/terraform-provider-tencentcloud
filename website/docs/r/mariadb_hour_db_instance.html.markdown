---
subcategory: "TencentDB for MariaDB(MariaDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mariadb_hour_db_instance"
sidebar_current: "docs-tencentcloud-resource-mariadb_hour_db_instance"
description: |-
  Provides a resource to create a MariaDB hour_db_instance
---

# tencentcloud_mariadb_hour_db_instance

Provides a resource to create a MariaDB hour_db_instance

## Example Usage

```hcl
resource "tencentcloud_mariadb_hour_db_instance" "basic" {
  db_version_id = "10.0"
  instance_name = "db-test-del"
  memory        = 2
  node_count    = 2
  storage       = 10
  subnet_id     = "subnet-jdi5xn22"
  vpc_id        = "vpc-k1t8ickr"
  vip           = "10.0.0.197"
  zones = [
    "ap-guangzhou-6",
    "ap-guangzhou-7",
  ]
  tags = {
    createdBy = "terraform"
  }
}
```

### Create with custom init_params

```hcl
resource "tencentcloud_mariadb_hour_db_instance" "with_init_params" {
  db_version_id = "10.0"
  instance_name = "db-test-init"
  memory        = 2
  node_count    = 2
  storage       = 10
  subnet_id     = "subnet-jdi5xn22"
  vpc_id        = "vpc-k1t8ickr"
  zones = [
    "ap-guangzhou-6",
    "ap-guangzhou-7",
  ]

  init_params {
    param = "character_set_server"
    value = "utf8mb4"
  }
  init_params {
    param = "lower_case_table_names"
    value = "1"
  }
  init_params {
    param = "innodb_page_size"
    value = "16384"
  }
  init_params {
    param = "sync_mode"
    value = "2"
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
* `init_params` - (Optional, List, ForceNew) parameter list. This interface's optional values include: `character_set_server` (character set, required), `lower_case_table_names` (table name case sensitivity, required, 0 - sensitive; 1 - insensitive), `innodb_page_size` (innodb data page, default 16K), `sync_mode` (sync mode: 0 - async; 1 - strong sync; 2 - strong sync degradable, default is strong sync degradable).
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



## Import

mariadb hour_db_instance can be imported using the id, e.g.
```
$ terraform import tencentcloud_mariadb_hour_db_instance.hour_db_instance tdsql-kjqih9nn
```

