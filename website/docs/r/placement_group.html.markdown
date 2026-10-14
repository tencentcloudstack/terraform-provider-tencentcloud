---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_placement_group"
sidebar_current: "docs-tencentcloud-resource-placement_group"
description: |-
  Provide a resource to create a placement group.
---

# tencentcloud_placement_group

Provide a resource to create a placement group.

## Example Usage

```hcl
resource "tencentcloud_placement_group" "foo" {
  name     = "test"
  type     = "HOST"
  affinity = 2
  tags = {
    createBy = "terraform"
  }
}
```

### Create partition placement group

```hcl
resource "tencentcloud_placement_group" "bar" {
  name            = "test-partition"
  type            = "HOST"
  strategy        = "PARTITION"
  partition_count = 5
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Name of the placement group, 1-60 characters in length.
* `type` - (Required, String, ForceNew) Type of the placement group. Valid values: `HOST`, `SW` and `RACK`.
* `affinity` - (Optional, Int, ForceNew) Affinity of the placement group.Valid values: 1~10, default is 1.
* `partition_count` - (Optional, Int) Partition count of the placement group. Valid values: 2~30. Only valid when `strategy` is set to `PARTITION`.
* `strategy` - (Optional, String) Strategy of the placement group. Valid values: `SPREAD` and `PARTITION`. `SPREAD` is the default strategy. When strategy is `PARTITION`, `partition_count` must be set. This field cannot be modified after creation.
* `tags` - (Optional, Map) Tags of the placement group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time of the placement group.
* `current_num` - Number of hosts in the placement group.
* `cvm_quota_total` - Maximum number of hosts in the placement group.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateDisasterRecoverGroup](https://www.tencentcloud.com/document/product/213/33264) - This API is used to create a [spread placement group](https://intl.cloud.tencent.com/document/product/213/15486?from_cn_redirect=1). After you create one, you can specify it for an instance when you [create the instance]
* Read: [DescribeDisasterRecoverGroups](https://www.tencentcloud.com/document/product/213/33261) - This API is used to query the information on [spread placement groups](https://intl.cloud.tencent.com/document/product/213/15486?from_cn_redirect=1).
* Update: [ModifyDisasterRecoverGroupAttribute](https://www.tencentcloud.com/document/product/213/33247) - This API is used to modify the attributes of [spread placement groups](https://intl.cloud.tencent.com/document/product/213/15486?from_cn_redirect=1).
* Update: [ModifyResourceTags](https://www.tencentcloud.com/document/product/651/33012) - This API is used to modify all tags associated with a resource.
* Delete: [DeleteDisasterRecoverGroups](https://www.tencentcloud.com/document/product/213/33263) - This API is used to delete a [spread placement group](https://intl.cloud.tencent.com/document/product/213/15486?from_cn_redirect=1). Only empty placement groups can be deleted.

## Import

Placement group can be imported using the id, e.g.

```
$ terraform import tencentcloud_placement_group.foo ps-ilan8vjf
```

