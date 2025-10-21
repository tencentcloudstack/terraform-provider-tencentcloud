---
subcategory: "TencentCloud Managed Service for Grafana(TCMG)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_grafana_notification_channel"
sidebar_current: "docs-tencentcloud-resource-monitor_grafana_notification_channel"
description: |-
  Provides a resource to create a monitor grafanaNotificationChannel
---

# tencentcloud_monitor_grafana_notification_channel

Provides a resource to create a monitor grafanaNotificationChannel

## Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-6"
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

resource "tencentcloud_monitor_grafana_instance" "foo" {
  instance_name         = "test-grafana"
  vpc_id                = tencentcloud_vpc.vpc.id
  subnet_ids            = [tencentcloud_subnet.subnet.id]
  grafana_init_password = "1234567890"
  enable_internet       = false

  tags = {
    "createdBy" = "test"
  }
}

resource "tencentcloud_monitor_alarm_notice" "foo" {
  name            = "tf_alarm_notice"
  notice_type     = "ALL"
  notice_language = "zh-CN"

  user_notices {
    receiver_type            = "USER"
    start_time               = 0
    end_time                 = 1
    notice_way               = ["SMS", "EMAIL"]
    user_ids                 = [10001]
    group_ids                = []
    phone_order              = [10001]
    phone_circle_times       = 2
    phone_circle_interval    = 50
    phone_inner_interval     = 60
    need_phone_arrive_notice = 1
    phone_call_type          = "CIRCLE"
    weekday                  = [1, 2, 3, 4, 5, 6, 7]
  }

  url_notices {
    url        = "https://www.mytest.com/validate"
    end_time   = 0
    start_time = 1
    weekday    = [1, 2, 3, 4, 5, 6, 7]
  }
}

resource "tencentcloud_monitor_grafana_notification_channel" "grafanaNotificationChannel" {
  instance_id   = tencentcloud_monitor_grafana_instance.foo.id
  channel_name  = "tf-channel"
  org_id        = 1
  receivers     = [tencentcloud_monitor_alarm_notice.foo.amp_consumer_id]
  extra_org_ids = ["1"]
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) grafana instance id.
* `channel_name` - (Optional, String) channel name.
* `extra_org_ids` - (Optional, Set: [`String`]) extra grafana organization id list, default to 1 representing Main Org.
* `org_id` - (Optional, Int) Grafana organization which channel will be installed, default to 1 representing Main Org.
* `receivers` - (Optional, Set: [`String`]) cloud monitor notification template notice-id list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `channel_id` - plugin id.


