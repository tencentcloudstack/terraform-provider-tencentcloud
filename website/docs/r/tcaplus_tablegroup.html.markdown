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

### Create a tcaplusdb table group

```hcl
resource "tencentcloud_tcaplus_cluster" "example" {
  idl_type                 = "PROTO"
  cluster_name             = "tf_example_tcaplus_cluster"
  vpc_id                   = "vpc-i5yyodl9"
  subnet_id                = "subnet-hhi88a58"
  password                 = "Password@2026"
  old_password_expire_last = 3600
}

resource "tencentcloud_tcaplus_tablegroup" "example" {
  cluster_id      = tencentcloud_tcaplus_cluster.example.id
  tablegroup_name = "tf_example_group_name"
  resource_tags {
    tag_key   = "CreatedBy"
    tag_value = "Terraform"
  }
}
```

### Create a tcaplusdb table group with user-specified table group id

```hcl
resource "tencentcloud_tcaplus_cluster" "example" {
  idl_type                 = "PROTO"
  cluster_name             = "tf_example_tcaplus_cluster"
  vpc_id                   = "vpc-i5yyodl9"
  subnet_id                = "subnet-hhi88a58"
  password                 = "Password@2026"
  old_password_expire_last = 3600
}

resource "tencentcloud_tcaplus_tablegroup" "example" {
  cluster_id      = tencentcloud_tcaplus_cluster.example.id
  tablegroup_name = "tf_example_group_name"
  table_group_id  = "109"
  resource_tags {
    tag_key   = "CreatedBy"
    tag_value = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) ID of the TcaplusDB cluster to which the table group belongs.
* `tablegroup_name` - (Required, String) Table group name; may consist of Chinese characters, English letters, or numeric characters, with a maximum length of 32 characters.
* `resource_tags` - (Optional, Set) Set of table group tags.
* `table_group_id` - (Optional, String, ForceNew) ID of the TcaplusDB table group, can be user-specified (must be unique within the cluster) or auto-incremented by the API when not set. Immutable after creation.

The `resource_tags` object supports the following:

* `tag_key` - (Required, String) Tag key.
* `tag_value` - (Required, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time of the TcaplusDB table group.
* `table_count` - Number of tables.
* `total_size` - Total storage size (MB).


## Import

TcaplusDB table group can be imported using the clusterId:tableGroupId, e.g.

```
terraform import tencentcloud_tcaplus_tablegroup.example 5516511420:52
```

