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

* Read:
  - [DescribeClusters](https://www.tencentcloud.com/search/?wd=DescribeClusters&type=document&lang=en)
  - [DescribeResourceTagsByResourceIds](https://www.tencentcloud.com/document/api/651/33022)
* Update:
  - [ModifyCluster](https://www.tencentcloud.com/search/?wd=ModifyCluster&type=document&lang=en)
  - [ModifyResourceTags](https://www.tencentcloud.com/document/api/651/33012)
* Delete: [DeleteCluster](https://www.tencentcloud.com/search/?wd=DeleteCluster&type=document&lang=en)

## Import

Tdmq instance can be imported, e.g.

```
$ terraform import tencentcloud_tdmq_instance.example pulsar-78bwjaj8epxv
```

