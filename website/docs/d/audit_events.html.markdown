---
subcategory: "Cloud Audit(Audit)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_audit_events"
sidebar_current: "docs-tencentcloud-datasource-audit_events"
description: |-
  Use this data source to query the events list supported by the audit.
---

# tencentcloud_audit_events

Use this data source to query the events list supported by the audit.

## Example Usage

```hcl
data "tencentcloud_audit_events" "events" {
  start_time  = "1727433841"
  end_time    = "1727437441"
  max_results = 50

  lookup_attributes {
    attribute_key   = "ResourceType"
    attribute_value = "cvm"
  }

  lookup_attributes {
    attribute_key   = "OnlyRecordNotSeen"
    attribute_value = "0"
  }

  lookup_attributes {
    attribute_key   = "EventPlatform"
    attribute_value = "0"
  }

  is_return_location = 1
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, Int) End timestamp in seconds (the time range for query is less than 30 days).
* `start_time` - (Required, Int) Start timestamp in seconds (cannot be 90 days after the current time).
* `is_return_location` - (Optional, Int) Whether to return the IP location. `1`: yes, `0`: no.
* `lookup_attributes` - (Optional, List) Search condition. Valid values: `RequestId`, `EventName`, `ActionType` (write/read), `PrincipalId` (sub-account), `ResourceType`, `ResourceName`, `AccessKeyId`, `SensitiveAction`, `ApiErrorCode`, `CamErrorCode`, and `Tags` (Format of AttributeValue: [{"key":"*","value":"*"}]).
* `max_results` - (Optional, Int) Max number of returned logs (up to 50).
* `result_output_file` - (Optional, String) Used to save results.

The `lookup_attributes` object supports the following:

* `attribute_key` - (Required, String) Valid values: RequestId, EventName, ReadOnly, Username, ResourceType, ResourceName, AccessKeyId, and EventId
Note: `null` may be returned for this field, indicating that no valid values can be obtained.
* `attribute_value` - (Optional, String) Value of `AttributeValue`
Note: `null` may be returned for this field, indicating that no valid values can be obtained.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `events` - Logset. Note: `null` may be returned for this field, indicating that no valid values can be obtained.


