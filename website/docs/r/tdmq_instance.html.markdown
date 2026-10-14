---
subcategory: "TDMQ for Pulsar(tpulsar)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_instance"
sidebar_current: "docs-tencentcloud-resource-tdmq_instance"
description: |-
  Provide a resource to create a TDMQ instance.
---

# tencentcloud_tdmq_instance

Provide a resource to create a TDMQ instance.

## Example Usage

```hcl
resource "tencentcloud_tdmq_instance" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
  tags = {
    "createdBy" = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_name` - (Required, String) The name of tdmq cluster to be created.
* `bind_cluster_id` - (Optional, Int) The Dedicated Cluster Id.
* `remark` - (Optional, String) Description of the tdmq cluster.
* `tags` - (Optional, Map) Tag description list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Read: [DescribeClusters](https://cloud.tencent.com/document/api/1179/52183) - This API is used to get the list of clusters.
* Read: [DescribeResourceTagsByResourceIds](https://www.tencentcloud.com/document/product/651/33022) - This API is used to query the tag key-value pairs associated with an existing resource.
* Update: [ModifyCluster](https://cloud.tencent.com/document/api/1179/52189) - This API is used to update a cluster.
* Update: [ModifyResourceTags](https://www.tencentcloud.com/document/product/651/33012) - This API is used to modify all tags associated with a resource.
* Delete: [DeleteCluster](https://cloud.tencent.com/document/api/1179/52185) - This API is used to delete a cluster.

## Import

Tdmq instance can be imported, e.g.

```
$ terraform import tencentcloud_tdmq_instance.example pulsar-78bwjaj8epxv
```

