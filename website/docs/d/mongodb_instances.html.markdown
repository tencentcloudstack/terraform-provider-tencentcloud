---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mongodb_instances"
sidebar_current: "docs-tencentcloud-datasource-mongodb_instances"
description: |-
  Use this data source to query detailed information of Mongodb instances.
---

# tencentcloud_mongodb_instances

Use this data source to query detailed information of Mongodb instances.

## Example Usage

```hcl
data "tencentcloud_mongodb_instances" "mongodb" {
  instance_id  = "cmgo-l6lwdsel"
  cluster_type = "REPLSET"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_type` - (Optional) Type of Mongodb cluster, and available values include replica set cluster(expressed with `REPLSET`), sharding cluster(expressed with `SHARD`).
* `instance_id` - (Optional) ID of the Mongodb instance to be queried.
* `instance_name_prefix` - (Optional) Name prefix of the Mongodb instance.
* `result_output_file` - (Optional) Used to store results.


