---
subcategory: "TDSQL for MySQL(DCDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dcdb_shard_spec"
sidebar_current: "docs-tencentcloud-datasource-dcdb_shard_spec"
description: |-
  Use this data source to query detailed information of dcdb shard_spec
---

# tencentcloud_dcdb_shard_spec

Use this data source to query detailed information of dcdb shard_spec

## Example Usage

```hcl
data "tencentcloud_dcdb_shard_spec" "shard_spec" {}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `spec_config` - list of instance specifications.
  * `machine` - machine type.
  * `spec_config_infos` - list of machine specifications.
    * `cpu` - CPU cores.
    * `max_storage` - maximum storage size, inGB.
    * `memory` - memory, in GB.
    * `min_storage` - minimum storage size, in GB.
    * `node_count` - node count.
    * `pid` - product price id.
    * `qps` - maximum QPS.
    * `suit_info` - recommended usage scenarios.


