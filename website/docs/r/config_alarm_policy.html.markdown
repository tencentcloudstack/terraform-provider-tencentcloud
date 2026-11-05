---
subcategory: "Config"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_config_alarm_policy"
sidebar_current: "docs-tencentcloud-resource-config_alarm_policy"
description: |-
  Provides a resource to create a Config alarm policy.
---

# tencentcloud_config_alarm_policy

Provides a resource to create a Config alarm policy.

## Example Usage

```hcl
resource "tencentcloud_config_alarm_policy" "example" {
  name                   = "tf-example"
  type                   = 1
  event_scope            = [1]
  risk_level             = [1, 2]
  notice_time            = "09:30:00~23:30:00"
  notification_mechanism = "Send in real time"
  status                 = 1
  notice_period          = [1, 2, 3, 4, 5]
  description            = "tf example alarm policy"
}
```

## Argument Reference

The following arguments are supported:

* `event_scope` - (Required, List: [`Int`]) Event scope. Valid values: 1 (current account), 2 (multi-account).
* `name` - (Required, String) Alarm policy name.
* `notice_period` - (Required, List: [`Int`]) Notification weekday list, 1-7 represent Monday to Sunday.
* `notice_time` - (Required, String) Notification time range, format: HH:mm:ss~HH:mm:ss (e.g. 09:30:00~23:30:00).
* `notification_mechanism` - (Required, String) Notification mechanism.
* `risk_level` - (Required, List: [`Int`]) Risk level list. Valid values: 1 (high risk), 2 (medium risk), 3 (low risk).
* `status` - (Required, Int) Policy status. Valid values: 1 (enabled), 2 (disabled).
* `type` - (Required, Int, ForceNew) Event Type 1: Resource Non Compliance Event.
* `description` - (Optional, String) Policy description.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `alarm_policy_id` - Alarm policy unique ID.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [AddAlarmPolicy](https://www.tencentcloud.com/document/api/1164/85014)
* Read: [ListAlarmPolicy](https://www.tencentcloud.com/document/api/1164/85013)
* Update: [UpdateAlarmPolicy](https://www.tencentcloud.com/document/api/1164/85012)
* Delete: [DeleteAlarmPolicy](https://www.tencentcloud.com/document/api/1164/84988)

## Import

Config alarm policy can be imported using the alarmPolicyId, e.g.

```
terraform import tencentcloud_config_alarm_policy.example 65
```

