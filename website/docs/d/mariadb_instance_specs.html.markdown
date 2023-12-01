---
subcategory: "TencentDB for MariaDB(MariaDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mariadb_instance_specs"
sidebar_current: "docs-tencentcloud-datasource-mariadb_instance_specs"
description: |-
  Use this data source to query detailed information of mariadb instance_specs
---

# tencentcloud_mariadb_instance_specs

Use this data source to query detailed information of mariadb instance_specs

## Example Usage

```hcl
data "tencentcloud_mariadb_instance_specs" "instance_specs" {
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `specs` - list of instance specifications.
  * `machine` - machine type.
  * `spec_infos` - list of machine specifications.
    * `cpu` - CPU cores.
    * `machine` - machine type.
    * `max_storage` - maximum storage size, in GB.
    * `memory` - memory, in GB.
    * `min_storage` - minimum storage size, in GB.
    * `node_count` - node count.
    * `pid` - product price id.
    * `qps` - maximum QPS.
    * `suit_info` - recommended usage scenarios.


