---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcaplus_group"
sidebar_current: "docs-tencentcloud-resource-tcaplus_group"
description: |-
  Use this resource to create tcaplus table group
---

# tencentcloud_tcaplus_group

Use this resource to create tcaplus table group

## Example Usage

```hcl
resource "tencentcloud_tcaplus_cluster" "test" {
  idl_type                 = "PROTO"
  cluster_name             = "tf_tcaplus_cluster_test"
  vpc_id                   = "vpc-7k6gzox6"
  subnet_id                = "subnet-akwgvfa3"
  password                 = "1qaA2k1wgvfa3ZZZ"
  old_password_expire_last = 3600
}

resource "tencentcloud_tcaplus_group" "group" {
  cluster_id = tencentcloud_tcaplus_cluster.test.id
  group_name = "tf_test_group_name"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, ForceNew) Cluster of the tcaplus group belongs.
* `group_name` - (Required) Name of the tcaplus group. length should between 1 and 30.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time of the tcaplus group.
* `table_count` - Number of tables.
* `total_size` - The total storage(MB).


