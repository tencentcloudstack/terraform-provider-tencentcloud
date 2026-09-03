Provides a resource to create a CLS scheduled sql

Example Usage

```hcl
resource "tencentcloud_cls_scheduled_sql" "example" {
  src_topic_id     = "2360e4a2-2b9e-4640-9829-72bc5052c9dd"
  src_topic_region = "ap-guangzhou"
  name             = "tf-example"
  enable_flag      = 1
  dst_resource {
    topic_id      = "3f5cc19d-f601-48b6-9671-9019b7b09bd0"
    region        = "ap-guangzhou"
    biz_type      = 1
    metric_names  = ["name1", "name2", "name3"]
    metric_labels = ["lable1", "lable2", "lable3"]
    custom_time   = "__WindowStartTime__"
    custom_metric_labels {
      key   = "cmlKey1"
      value = "cmlValue1"
    }

    custom_metric_labels {
      key   = "cmlKey2"
      value = "cmlValue2"
    }

    custom_metric_labels {
      key   = "cmlKey3"
      value = "cmlValue3"
    }
  }

  scheduled_sql_content = "verb:get AND responseStatus.code>=400\n| select stageTimestamp, responseStatus.code, objectRef.resource, objectRef.name, objectRef.namespace, user.username, \"annotations.authorization.k8s.io/decision\" as auth_decision\n  order by stageTimestamp desc\n  limit 50"
  process_start_time    = 1788417300000
  process_type          = 1
  process_period        = 1
  process_time_window   = "@m-1m,@m"
  process_delay         = 60
  syntax_rule           = 1
}
```

Import

CLS scheduled sql can be imported using the id, e.g.

```
terraform import tencentcloud_cls_scheduled_sql.example aebbad1b-4228-4ccc-8d62-0333c9739452
```