---
subcategory: "TencentDB for MariaDB(MariaDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mariadb_hour_db_instance"
sidebar_current: "docs-tencentcloud-resource-mariadb_hour_db_instance"
description: |-
  Provides a resource to create a mariadb hour_db_instance
---

# tencentcloud_mariadb_hour_db_instance

Provides a resource to create a mariadb hour_db_instance

## Example Usage

```hcl
resource "tencentcloud_mariadb_hour_db_instance" "basic" {
  db_version_id = "8.0"
  instance_name = "db-test-2"
  memory        = 2
  node_count    = 2
  storage       = 10
  subnet_id     = "subnet-jdi5xn22"
  vpc_id        = "vpc-k1t8ickr"
  zones         = ["ap-guangzhou-7", "ap-guangzhou-7"]
  tags = {
    createdBy = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `memory` - (Required, Int) instance memory.
* `node_count` - (Required, Int) number of node for instance.
* `storage` - (Required, Int) instance disk storage.
* `zones` - (Required, Set: [`String`]) available zone of instance.
* `db_version_id` - (Optional, String) db engine version, default to 10.1.9.
* `instance_name` - (Optional, String) name of this instance.
* `subnet_id` - (Optional, String) subnet id, it&amp;#39;s required when vpcId is set.
* `tags` - (Optional, Map) Tag description list.
* `vpc_id` - (Optional, String) vpc id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mariadb hour_db_instance can be imported using the id, e.g.
```
$ terraform import tencentcloud_mariadb_hour_db_instance.hour_db_instance tdsql-kjqih9nn
```

