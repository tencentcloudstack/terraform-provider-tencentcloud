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
* `subnet_id` - (Optional, String) subnet id, it&amp;#39;s required when vpcId is set.
* `tags` - (Optional, Map) Tag description list.
* `vpc_id` - (Optional, String) vpc id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mariadb dedicatedcluster_db_instance can be imported using the id, e.g.
```
$ terraform import tencentcloud_mariadb_dedicatedcluster_db_instance.dedicatedcluster_db_instance tdsql-050g3fmv
```

