Provides a resource to create a cls config extra

Example Usage

```hcl
resource "tencentcloud_cls_logset" "logset" {
  logset_name = "tf-config-extra-test"
  tags        = {
    "test" = "test"
  }
}

resource "tencentcloud_cls_topic" "topic" {
  auto_split           = true
  logset_id            = tencentcloud_cls_logset.logset.id
  max_split_partitions = 20
  partition_count      = 1
  period               = 10
  storage_type         = "hot"
  tags                 = {
    "test" = "test"
  }
  topic_name = "tf-config-extra-test"
}

resource "tencentcloud_cls_machine_group" "group" {
  group_name        = "tf-config-extra-test"
  service_logging   = true
  auto_update       = true
  update_end_time   = "19:05:00"
  update_start_time = "17:05:00"

  machine_group_type {
    type   = "ip"
    values = [
      "192.168.1.1",
      "192.168.1.2",
    ]
  }
}

resource "tencentcloud_cls_config_extra" "extra" {
  name        = "helloworld-test"
  topic_id    = tencentcloud_cls_topic.topic.id
  type        = "container_file"
  log_type    = "json_log"
  config_flag = "label_k8s"
  logset_id   = tencentcloud_cls_logset.logset.id
  logset_name = tencentcloud_cls_logset.logset.logset_name
  topic_name  = tencentcloud_cls_topic.topic.topic_name
  container_file {
    container    = "nginx"
    file_pattern = "log"
    log_path     = "/nginx"
    namespace    = "default"
    workload {
      container = "nginx"
      kind      = "deployment"
      name      = "nginx"
      namespace = "default"
    }
  }
  group_id = tencentcloud_cls_machine_group.group.id
}
```

Import

cls config_extra can be imported using the id, e.g.

```
terraform import tencentcloud_cls_config_extra.config_extra config_extra_id
```