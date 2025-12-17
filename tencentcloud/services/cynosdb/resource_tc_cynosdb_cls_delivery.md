Provides a resource to create a CynosDB cls delivery

~> **NOTE:** After executing `terraform destroy`, slow logs will no longer be uploaded, but historical logs will continue to be stored in the log topic until they expire. Log storage fees will continue to be charged during this period. If you do not wish to continue storing historical logs, you can go to CLS to delete the log topic.

Example Usage

Use topic_name and group_name

```hcl
resource "tencentcloud_cynosdb_cls_delivery" "example" {
  instance_id = "cynosdbmysql-ins-m2903cxq"
  cls_info_list {
    region     = "ap-guangzhou"
    topic_name = "tf-example"
    group_name = "tf-example"
  }

  running_status = true
}
```

Use topic_id and group_id

```hcl
resource "tencentcloud_cynosdb_cls_delivery" "example" {
  instance_id = "cynosdbmysql-ins-m2903cxq"
  cls_info_list {
    region   = "ap-guangzhou"
    topic_id = "a9d582f8-8c14-462c-94b8-bbc579a04f02"
    group_id = "67fca013-379b-4bc6-8e72-390227d869c4"
  }

  running_status = false
}
```

Import

CynosDB cls delivery can be imported using the instanceId#topicId, e.g.

```
terraform import tencentcloud_cynosdb_cls_delivery.example cynosdbmysql-ins-m2903cxq#222932ff-a10a-41f1-8d29-ff0cfe2a2d99
```
