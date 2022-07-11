---
subcategory: "TcaplusDB"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcaplus_tablegroup"
sidebar_current: "docs-tencentcloud-resource-tcaplus_tablegroup"
description: |-
  Use this resource to create TcaplusDB table group.
---

# tencentcloud_tcaplus_tablegroup

Use this resource to create TcaplusDB table group.

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

resource "tencentcloud_tcaplus_tablegroup" "tablegroup" {
  cluster_id      = tencentcloud_tcaplus_cluster.test.id
  tablegroup_name = "tf_test_group_name"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) ID of the TcaplusDB cluster to which the table group belongs.
* `tablegroup_name` - (Required, String) Name of the TcaplusDB table group. Name length should be between 1 and 30.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time of the TcaplusDB table group.
* `table_count` - Number of tables.
* `total_size` - Total storage size (MB).


