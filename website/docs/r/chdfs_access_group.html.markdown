---
subcategory: "Cloud HDFS(CHDFS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_chdfs_access_group"
sidebar_current: "docs-tencentcloud-resource-chdfs_access_group"
description: |-
  Provides a resource to create a chdfs access_group
---

# tencentcloud_chdfs_access_group

Provides a resource to create a chdfs access_group

## Example Usage

```hcl
resource "tencentcloud_chdfs_access_group" "access_group" {
  access_group_name = "testAccessGroup"
  vpc_type          = 1
  vpc_id            = "vpc-4owdpnwr"
  description       = "test access group"
}
```

## Argument Reference

The following arguments are supported:

* `access_group_name` - (Required, String) Permission group name.
* `vpc_id` - (Required, String) VPC ID.
* `vpc_type` - (Required, Int) vpc network type(1:CVM, 2:BM 1.0).
* `description` - (Optional, String) Permission group description, default empty.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateAccessGroup](https://cloud.tencent.com/document/api/1105/51143) - This API is used to create a permission group.
* Read: [DescribeAccessGroup](https://cloud.tencent.com/document/api/1105/51142) - This API is used to view the details of a permission group.
* Update: [ModifyAccessGroup](https://cloud.tencent.com/document/api/1105/51140) - This API is used to modify the attributes of a permission group.
* Delete: [DeleteAccessGroup](https://cloud.tencent.com/document/api/1105/52340) - This API is used to delete a permission group.

## Import

chdfs access_group can be imported using the id, e.g.

```
terraform import tencentcloud_chdfs_access_group.access_group access_group_id
```

