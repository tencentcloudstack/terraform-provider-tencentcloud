---
subcategory: "Cloud HDFS(CHDFS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_chdfs_access_groups"
sidebar_current: "docs-tencentcloud-datasource-chdfs_access_groups"
description: |-
  Use this data source to query detailed information of chdfs access_groups
---

# tencentcloud_chdfs_access_groups

Use this data source to query detailed information of chdfs access_groups

## Example Usage

```hcl
data "tencentcloud_chdfs_access_groups" "access_groups" {
  vpc_id = "vpc-pewdpc0d"
}
```

## Argument Reference

The following arguments are supported:

* `owner_uin` - (Optional, Int) get groups belongs to the owner uin, must set but only can use one of VpcId and OwnerUin to get the groups.
* `result_output_file` - (Optional, String) Used to save results.
* `vpc_id` - (Optional, String) get groups belongs to the vpc id, must set but only can use one of VpcId and OwnerUin to get the groups.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `access_groups` - access group list.
  * `access_group_id` - access group id.
  * `access_group_name` - access group name.
  * `create_time` - create time.
  * `description` - access group description.
  * `vpc_id` - VPC ID.
  * `vpc_type` - vpc network type(1:CVM, 2:BM 1.0).


