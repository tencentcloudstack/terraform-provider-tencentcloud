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

* Read: [DescribeClusters](https://cloud.tencent.com/document/api/1179/52183)
* Read: [DescribeResourceTagsByResourceIds](https://cloud.tencent.com/document/api/651/35319)
* Update: [ModifyCluster](https://cloud.tencent.com/document/api/1179/52189)
* Update: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322)
* Delete: [DeleteCluster](https://cloud.tencent.com/document/api/1179/52185)

## Import

Tdmq instance can be imported, e.g.

```
$ terraform import tencentcloud_tdmq_instance.example pulsar-78bwjaj8epxv
```

