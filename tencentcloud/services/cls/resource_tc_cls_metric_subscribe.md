Provides a resource to create a CLS metric subscribe

Example Usage

```hcl
resource "tencentcloud_cls_metric_subscribe" "example" {
  name      = "tf-example-metric-subscribe"
  topic_id  = "c9b68233-948a-4eaf-a363-d0c2ced393ae"
  namespace = "QCE/CVM"
  enable    = 2
  metrics {
    metric_name = "cpu_usage"
    periods     = [60, 300]
    metric_labels {
      key   = "label_key"
      value = "label_value"
    }
  }
  instance_info {
    instance_dimension = ["InstanceId"]
    instances {
      values = ["ins-xxxxxxxx"]
    }
  }
}
```

Import

CLS metric subscribe can be imported using the composite id, e.g. the composite id is `topicId#taskId` formatted by `topicId` and `taskId` joined with `#`:

```
terraform import tencentcloud_cls_metric_subscribe.example topicId#taskId
```
