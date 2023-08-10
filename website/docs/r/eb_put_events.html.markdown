---
subcategory: "EventBridge(EB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eb_put_events"
sidebar_current: "docs-tencentcloud-resource-eb_put_events"
description: |-
  Provides a resource to create a eb put_events
---

# tencentcloud_eb_put_events

Provides a resource to create a eb put_events

## Example Usage

```hcl
resource "tencentcloud_eb_event_bus" "foo" {
  event_bus_name = "tf-event_bus"
  description    = "event bus desc"
  enable_store   = false
  save_days      = 1
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_eb_put_events" "put_events" {
  event_list {
    source = "ckafka.cloud.tencent"
    data = jsonencode(
      {
        "topic" : "test-topic",
        "Partition" : 1,
        "offset" : 37,
        "msgKey" : "test",
        "msgBody" : "Hello from Ckafka again!"
      }
    )
    type    = "connector:ckafka"
    subject = "qcs::ckafka:ap-guangzhou:uin/1250000000:ckafkaId/uin/1250000000/ckafka-123456"
    time    = 1691572461939

  }
  event_bus_id = tencentcloud_eb_event_bus.foo.id
}
```

## Argument Reference

The following arguments are supported:

* `event_bus_id` - (Required, String, ForceNew) event bus Id.
* `event_list` - (Required, List, ForceNew) event list.

The `event_list` object supports the following:

* `data` - (Required, String) Event data, the content is controlled by the system that created the event, the current datacontenttype only supports application/json;charset=utf-8, so this field is a json string.
* `source` - (Required, String) Event source information, new product reporting must comply with EB specifications.
* `subject` - (Required, String) Detailed description of the event source, customizable, optional. The cloud service defaults to the standard qcs resource representation syntax: qcs::dts:ap-guangzhou:appid/uin:xxx.
* `type` - (Required, String) Event type, customizable, optional. The cloud service writes COS:Created:PostObject by default, use: to separate the type field.
* `time` - (Optional, Int) The timestamp in milliseconds when the event occurred,time.Now().UnixNano()/1e6.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



