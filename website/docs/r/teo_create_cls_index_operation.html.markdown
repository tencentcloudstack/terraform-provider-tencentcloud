---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_create_cls_index_operation"
sidebar_current: "docs-tencentcloud-resource-teo_create_cls_index_operation"
description: |-
  Provides a resource to create CLS index for TEO realtime log delivery task.
---

# tencentcloud_teo_create_cls_index_operation

Provides a resource to create CLS index for TEO realtime log delivery task.

## Example Usage

```hcl
resource "tencentcloud_teo_realtime_log_delivery" "example" {
  area            = "overseas"
  delivery_status = "enabled"
  entity_list = [
    "sid-2yvhjw98uaco",
  ]
  fields = [
    "ServiceID",
    "ConnectTimeStamp",
    "DisconnetTimeStamp",
    "DisconnetReason",
    "ClientRealIP",
    "ClientRegion",
    "EdgeIP",
    "ForwardProtocol",
    "ForwardPort",
    "SentBytes",
    "ReceivedBytes",
    "LogTimeStamp",
  ]
  log_type  = "application"
  sample    = 0
  task_name = "test-task"
  task_type = "cls"
  zone_id   = "zone-2qtuhspy7cr6"

  log_format {
    type      = "json"
    delimiter = ""
  }

  cls {
    region           = "ap-guangzhou"
    log_set_id       = "xxxxxxxxxx"
    topic_id         = "xxxxxxxxxx"
    enable_timestamp = true
  }
}

resource "tencentcloud_teo_create_cls_index_operation" "example" {
  zone_id = tencentcloud_teo_realtime_log_delivery.example.zone_id
  task_id = tencentcloud_teo_realtime_log_delivery.example.task_id
}
```

## Argument Reference

The following arguments are supported:

* `task_id` - (Required, String, ForceNew) Realtime log delivery task ID.
* `zone_id` - (Required, String, ForceNew) Zone ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

The resource can be imported by using the `zone_id`, e.g.

```sh
terraform import tencentcloud_teo_create_cls_index_operation.example zone-12345678
```

