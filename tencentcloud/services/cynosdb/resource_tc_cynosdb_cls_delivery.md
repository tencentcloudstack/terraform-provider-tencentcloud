Provides a resource to create a Cynosdb cls delivery

~> **NOTE:** Destroying this resource will not cause the CLS log set and log topic to be destroyed synchronously. If you need to delete it, you need to access the console page to delete it.

Example Usage

If topic_operation/group_operation is create

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

If topic_operation/group_operation is reuse

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

Import

Cynosdb cls delivery can be imported using the ${instanceId}#${groupId}#${topicId}, e.g.

```
terraform import tencentcloud_cynosdb_cls_delivery.example cynosdbmysql-ins-anknkhvi#7e3bb8b7-81d5-4e6b-b150-f139b90c146e#8e38f7c1-17ec-4acb-a4cb-7dc14baaef47
```
