Provides a resource to create a eb event_target

Example Usage

Create an event target of type scf

```hcl
variable "zone" {
  default = "ap-guangzhou"
}

variable "namespace" {
  default = "default"
}

variable "function" {
  default = "keep-1676351130"
}

variable "function_version" {
  default = "$LATEST"
}

data "tencentcloud_cam_users" "foo" {
}

resource "tencentcloud_eb_event_bus" "foo" {
  event_bus_name = "tf-event_bus"
  description    = "event bus desc"
  enable_store   = false
  save_days      = 1
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_eb_event_rule" "foo" {
  event_bus_id = tencentcloud_eb_event_bus.foo.id
  rule_name    = "tf-event_rule"
  description  = "event rule desc"
  enable       = true
  event_pattern = jsonencode(
    {
      source = "apigw.cloud.tencent"
      type = [
        "connector:apigw",
      ]
    }
  )
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_eb_event_target" "scf_target" {
    event_bus_id = tencentcloud_eb_event_bus.foo.id
    rule_id      = tencentcloud_eb_event_rule.foo.rule_id
    type         = "scf"

    target_description {
        resource_description = "qcs::scf:${var.zone}:uin/${data.tencentcloud_cam_users.foo.user_list.0.uin}:namespace/${var.namespace}/function/${var.function}/${var.function_version}"

        scf_params {
            batch_event_count     = 1
            batch_timeout         = 1
            enable_batch_delivery = true
        }
    }
}
```

Create an event target of type ckafka

```hcl
variable "ckafka" {
  default = "ckafka-qzoeaqx8"
}

resource "tencentcloud_eb_event_target" "ckafka_target" {
    event_bus_id = tencentcloud_eb_event_bus.foo.id
    rule_id      = tencentcloud_eb_event_rule.foo.rule_id
    type         = "ckafka"

    target_description {
        resource_description = "qcs::scf:${var.zone}:uin/${data.tencentcloud_cam_users.foo.user_list.0.uin}:ckafkaId/uin/${data.tencentcloud_cam_users.foo.user_list.0.uin}/${var.ckafka}"

        ckafka_target_params {
            topic_name = "dasdasd"

            retry_policy {
                max_retry_attempts = 360
                retry_interval     = 60
            }
        }
    }
}
```

Import

eb event_target can be imported using the id, e.g.

```
terraform import tencentcloud_eb_event_target.event_target event_target_id
```