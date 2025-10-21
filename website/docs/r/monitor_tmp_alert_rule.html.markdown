---
subcategory: "Managed Service for Prometheus(TMP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_tmp_alert_rule"
sidebar_current: "docs-tencentcloud-resource-monitor_tmp_alert_rule"
description: |-
  Provides a resource to create a monitor tmpAlertRule
---

# tencentcloud_monitor_tmp_alert_rule

Provides a resource to create a monitor tmpAlertRule

~> **NOTE:** This resource will be deprecated in a future version, please use `tencentcloud_monitor_tmp_alert_group`.

## Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-4"
}

resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_monitor_vpc"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  name              = "tf_monitor_subnet"
  cidr_block        = "10.0.1.0/24"
}

resource "tencentcloud_monitor_tmp_instance" "foo" {
  instance_name       = "tf-tmp-instance"
  vpc_id              = tencentcloud_vpc.vpc.id
  subnet_id           = tencentcloud_subnet.subnet.id
  data_retention_time = 30
  zone                = var.availability_zone
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_monitor_tmp_cvm_agent" "foo" {
  instance_id = tencentcloud_monitor_tmp_instance.foo.id
  name        = "tf-agent"
}

resource "tencentcloud_monitor_tmp_alert_rule" "foo" {
  duration    = "2m"
  expr        = "avg by (instance) (mysql_global_status_threads_connected) / avg by (instance) (mysql_global_variables_max_connections)  > 0.8"
  instance_id = tencentcloud_monitor_tmp_instance.foo.id
  receivers   = ["notice-f2svbu3w"]
  rule_name   = "MySQL 连接数过多"
  rule_state  = 2
  type        = "MySQL/MySQL 连接数过多"

  annotations {
    key   = "description"
    value = "MySQL 连接数过多, 实例: {{$labels.instance}}，当前值: {{ $value | humanizePercentage }}。"
  }
  annotations {
    key   = "summary"
    value = "MySQL 连接数过多(>80%)"
  }

  labels {
    key   = "severity"
    value = "warning"
  }
}
```

## Argument Reference

The following arguments are supported:

* `expr` - (Required, String) Rule expression, reference documentation: `https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules/`.
* `instance_id` - (Required, String) Instance id.
* `receivers` - (Required, Set: [`String`]) Alarm notification template id list.
* `rule_name` - (Required, String) Rule name.
* `annotations` - (Optional, Set) Rule alarm duration.
* `duration` - (Optional, String) Rule alarm duration.
* `labels` - (Optional, Set) Rule alarm duration.
* `rule_state` - (Optional, Int) Rule state code.
* `type` - (Optional, String) Alarm Policy Template Classification.

The `annotations` object supports the following:

* `key` - (Required, String) key.
* `value` - (Required, String) value.

The `labels` object supports the following:

* `key` - (Required, String) key.
* `value` - (Required, String) value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

monitor tmpAlertRule can be imported using the id, e.g.
```
$ terraform import tencentcloud_monitor_tmp_alert_rule.tmpAlertRule instanceId#Rule_id
```

