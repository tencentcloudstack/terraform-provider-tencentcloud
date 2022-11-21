---
subcategory: "TDSQL for MySQL(dcdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dcdb_hourdb_instance"
sidebar_current: "docs-tencentcloud-resource-dcdb_hourdb_instance"
description: |-
  Provides a resource to create a dcdb hourdb_instance
---

# tencentcloud_dcdb_hourdb_instance

Provides a resource to create a dcdb hourdb_instance

## Example Usage

```hcl
resource "tencentcloud_dcdb_hourdb_instance" "hourdb_instance" {
  instance_name     = "test_dcdc_dc_instance"
  zones             = ["ap-guangzhou-5", "ap-guangzhou-6"]
  shard_memory      = "2"
  shard_storage     = "10"
  shard_node_count  = "2"
  shard_count       = "2"
  vpc_id            = local.vpc_id
  subnet_id         = local.subnet_id
  security_group_id = local.sg_id
  db_version_id     = "8.0"
  resource_tags {
    tag_key   = "aaa"
    tag_value = "bbb"
  }
}
```

## Argument Reference

The following arguments are supported:

* `shard_count` - (Required, Int) instance shard count.
* `shard_memory` - (Required, Int) memory(GB) for each shard. It can be obtained by querying api DescribeShardSpec.
* `shard_node_count` - (Required, Int) node count for each shard. It can be obtained by querying api DescribeShardSpec.
* `shard_storage` - (Required, Int) storage(GB) for each shard. It can be obtained by querying api DescribeShardSpec.
* `db_version_id` - (Optional, String) db engine version, default to Percona 5.7.17.
* `instance_name` - (Optional, String) name of this instance.
* `project_id` - (Optional, Int) project id.
* `resource_tags` - (Optional, List) resource tags.
* `security_group_id` - (Optional, String) security group id.
* `subnet_id` - (Optional, String) subnet id, it&amp;#39;s required when vpcId is set.
* `vpc_id` - (Optional, String) vpc id.
* `zones` - (Optional, Set: [`String`]) available zone.

The `resource_tags` object supports the following:

* `tag_key` - (Required, String) tag key.
* `tag_value` - (Required, String) tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

dcdb hourdb_instance can be imported using the id, e.g.
```
$ terraform import tencentcloud_dcdb_hourdb_instance.hourdb_instance hourdbInstance_id
```

