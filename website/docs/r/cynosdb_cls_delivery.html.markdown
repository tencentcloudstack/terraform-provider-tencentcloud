---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_cls_delivery"
sidebar_current: "docs-tencentcloud-resource-cynosdb_cls_delivery"
description: |-
  Provides a resource to create a CynosDB cls delivery
---

# tencentcloud_cynosdb_cls_delivery

Provides a resource to create a CynosDB cls delivery

~> **NOTE:** After executing `terraform destroy`, slow logs will no longer be uploaded, but historical logs will continue to be stored in the log topic until they expire. Log storage fees will continue to be charged during this period. If you do not wish to continue storing historical logs, you can go to CLS to delete the log topic.

## Example Usage

### Use topic_name and group_name

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

### Use topic_id and group_id

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

## Argument Reference

The following arguments are supported:

* `cls_info_list` - (Required, List, ForceNew) Log shipping configuration.
* `instance_id` - (Required, String, ForceNew) Intance ID.
* `log_type` - (Optional, String, ForceNew) Log type.
* `running_status` - (Optional, Bool) Delivery status. true: Enabled; false: Disabled.

The `cls_info_list` object supports the following:

* `region` - (Required, String, ForceNew) Log delivery area.
* `group_id` - (Optional, String, ForceNew) Log set ID.
* `group_name` - (Optional, String, ForceNew) Log set name.
* `topic_id` - (Optional, String, ForceNew) Log topic ID.
* `topic_name` - (Optional, String, ForceNew) Log topic name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

CynosDB cls delivery can be imported using the instanceId#topicId, e.g.

```
terraform import tencentcloud_cynosdb_cls_delivery.example cynosdbmysql-ins-m2903cxq#222932ff-a10a-41f1-8d29-ff0cfe2a2d99
```

