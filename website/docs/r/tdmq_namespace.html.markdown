---
subcategory: "TDMQ for Pulsar(tpulsar)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_namespace"
sidebar_current: "docs-tencentcloud-resource-tdmq_namespace"
description: |-
  Provide a resource to create a TDMQ namespace.
---

# tencentcloud_tdmq_namespace

Provide a resource to create a TDMQ namespace.

## Example Usage

```hcl
resource "tencentcloud_tdmq_instance" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tdmq_namespace" "example" {
  environ_name = "tf_example"
  msg_ttl      = 300
  cluster_id   = tencentcloud_tdmq_instance.example.id
  retention_policy {
    time_in_minutes = 60
    size_in_mb      = 10
  }
  remark = "remark."
  tags {
    tag_key   = "createdBy"
    tag_value = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) The Dedicated Cluster Id.
* `environ_name` - (Required, String) The name of namespace to be created.
* `msg_ttl` - (Required, Int) The expiration time of unconsumed message.
* `remark` - (Optional, String) Description of the namespace.
* `retention_policy` - (Optional, List) The Policy of message to retain. Format like: `{time_in_minutes: Int, size_in_mb: Int}`. `time_in_minutes`: the time of message to retain; `size_in_mb`: the size of message to retain.
* `tags` - (Optional, List) The tags of the tencentcloud_tdmq_namespace.

The `retention_policy` object supports the following:

* `size_in_mb` - (Optional, Int) the size of message to retain.
* `time_in_minutes` - (Optional, Int) the time of message to retain.

The `tags` object supports the following:

* `tag_key` - (Required, String) Tag key.
* `tag_value` - (Required, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateEnvironment](https://cloud.tencent.com/document/api/1179/46081) - This API is used to create a TDMQ namespace.
* Read: [DescribeEnvironments](https://cloud.tencent.com/document/api/1179/46078) - This API is used to get the list of namespaces under a tenant.
* Update: [ModifyEnvironmentAttributes](https://cloud.tencent.com/document/api/1179/46077) - This API is used to modify the attributes of the specified namespace.
* Update: [ModifyResourceTags](https://www.tencentcloud.com/document/product/651/33012) - This API is used to modify all tags associated with a resource.
* Delete: [DeleteEnvironments](https://cloud.tencent.com/document/api/1179/46080) - This API is used to batch delete namespaces under a tenant.

## Import

Tdmq namespace can be imported, e.g.

```
$ terraform import tencentcloud_tdmq_instance.example tf_example#pulsar-78bwjaj8epxv
```

