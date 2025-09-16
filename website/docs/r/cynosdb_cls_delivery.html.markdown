---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_cls_delivery"
sidebar_current: "docs-tencentcloud-resource-cynosdb_cls_delivery"
description: |-
  Provides a resource to create a Cynosdb cls delivery
---

# tencentcloud_cynosdb_cls_delivery

Provides a resource to create a Cynosdb cls delivery

~> **NOTE:** Destroying this resource will not cause the CLS log set and log topic to be destroyed synchronously. If you need to delete it, you need to access the console page to delete it.

## Example Usage

### If topic_operation/group_operation is create

```hcl
resource "tencentcloud_cynosdb_cls_delivery" "example" {
  instance_id = "cynosdbmysql-ins-anknkhvi"
  cls_info_list {
    topic_operation = "create"
    group_operation = "create"
    region          = "ap-guangzhou"
    topic_name      = "tf-example-topic"
    group_name      = "tf-example-group"
  }
  log_type            = "slow"
  enable_cls_delivery = true
}
```

### If topic_operation/group_operation is reuse

```hcl
resource "tencentcloud_cynosdb_cls_delivery" "example" {
  instance_id = "cynosdbmysql-ins-anknkhvi"
  cls_info_list {
    topic_operation = "reuse"
    group_operation = "reuse"
    region          = "ap-guangzhou"
    topic_id        = "8e38f7c1-17ec-4acb-a4cb-7dc14baaef47"
    group_id        = "7e3bb8b7-81d5-4e6b-b150-f139b90c146e"
  }
  log_type            = "slow"
  enable_cls_delivery = false
}
```

## Argument Reference

The following arguments are supported:

* `cls_info_list` - (Required, List, ForceNew) Log delivery configuration.
* `instance_id` - (Required, String, ForceNew) Instance ID.
* `enable_cls_delivery` - (Optional, Bool) Whether to enable CLS delivery. Default value: true (enabled).
* `log_type` - (Optional, String, ForceNew) Log type.

The `cls_info_list` object supports the following:

* `group_operation` - (Required, String) Log set operations: Optional: create or reuse. create: Creates a new log set, using the GroupName. reuse: Reuses an existing log set, using the GroupId. The combination of reusing an existing log topic and creating a new log set is not allowed.
* `region` - (Required, String) Log delivery region.
* `topic_operation` - (Required, String) Log topic operations: Optional: create or reuse. create: Creates a new log topic using TopicName. reuse: Reuses an existing log topic using TopicId. Reusing an existing log topic and creating a new log set is not allowed.
* `group_id` - (Optional, String) Log set ID.
* `group_name` - (Optional, String) Log set name.
* `topic_id` - (Optional, String) Log topic ID.
* `topic_name` - (Optional, String) Log topic name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `status` - Delivery status. running, offlined.


## Import

Cynosdb cls delivery can be imported using the ${instanceId}#${groupId}#${topicId}, e.g.

```
terraform import tencentcloud_cynosdb_cls_delivery.example cynosdbmysql-ins-anknkhvi#7e3bb8b7-81d5-4e6b-b150-f139b90c146e#8e38f7c1-17ec-4acb-a4cb-7dc14baaef47
```

