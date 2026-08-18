---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_metric_subscribe"
sidebar_current: "docs-tencentcloud-resource-cls_metric_subscribe"
description: |-
  Provides a resource to create a CLS metric subscribe
---

# tencentcloud_cls_metric_subscribe

Provides a resource to create a CLS metric subscribe

## Example Usage

```hcl
resource "tencentcloud_cls_metric_subscribe" "example" {
  name      = "tf-example"
  topic_id  = "da978f4a-b205-4286-830c-86d0ca45d7f3"
  namespace = "qce/cvm"
  enable    = 2
  metrics {
    metric_name = "BaseCpuUsage"
    periods     = [10]
    metric_labels {
      key   = "label_key"
      value = "label_value"
    }
  }

  metrics {
    metric_name = "CbsVolumeFsUsage"
    periods     = [10]
    metric_labels {
      key   = "label_key"
      value = "label_value"
    }
  }

  metrics {
    metric_name = "MemUsed"
    periods     = [10]
    metric_labels {
      key   = "label_key"
      value = "label_value"
    }
  }

  instance_info {
    instance_dimension = [
      "InstanceId",
      "diskid"
    ]

    instances {
      values = [
        "ins-ly9ibb5w",
        "disk-781trig2"
      ]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_info` - (Required, List) Instance config info.
* `metrics` - (Required, List) Metric config info list.
* `name` - (Required, String) Subscribe task name, up to 64 characters, start with a letter, support 0-9, a-z, A-Z, _, -, Chinese characters.
* `namespace` - (Required, String) Cloud product namespace.
* `topic_id` - (Required, String, ForceNew) Log topic id.
* `enable` - (Optional, Int) Task switch, 1: pause, 2: enable.

The `instance_info` object supports the following:

* `instance_dimension` - (Required, List) Instance dimension.
* `instances` - (Required, List) Instance value list.

The `instances` object of `instance_info` supports the following:

* `values` - (Required, List) Instance info value list.

The `metric_labels` object of `metrics` supports the following:

* `key` - (Required, String) Metric label name.
* `value` - (Required, String) Metric label content.

The `metrics` object supports the following:

* `metric_name` - (Required, String) Metric name.
* `periods` - (Required, List) Statistical period, unit: second(s).
* `metric_labels` - (Optional, List) Custom metric labels.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `status` - Subscribe task running status. 0: creating, 1: paused, 2: running, 3: abnormal.
* `task_id` - Subscribe task id.


## Import

CLS metric subscribe can be imported using the composite id, e.g. the composite id is `topicId#taskId` formatted by `topicId` and `taskId` joined with `#`:

```
terraform import tencentcloud_cls_metric_subscribe.example da978f4a-b205-4286-830c-86d0ca45d7f3#bc28af29-b274-4a69-ba77-6e665db26850
```

