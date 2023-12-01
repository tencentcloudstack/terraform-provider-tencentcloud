---
subcategory: "EventBridge(EB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eb_bus"
sidebar_current: "docs-tencentcloud-datasource-eb_bus"
description: |-
  Use this data source to query detailed information of eb bus
---

# tencentcloud_eb_bus

Use this data source to query detailed information of eb bus

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
data "tencentcloud_eb_bus" "bus" {
  order_by = "AddTime"
  order    = "DESC"
  filters {
    values = ["Custom"]
    name   = "Type"
  }

  depends_on = [tencentcloud_eb_event_bus.foo]
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter conditions. The upper limit of Filters per request is 10, and the upper limit of Filter.Values 5.
* `order_by` - (Optional, String) According to which field to sort the returned results, the following fields are supported: AddTime (creation time), ModTime (modification time).
* `order` - (Optional, String) Return results in ascending or descending order, optional values ASC (ascending) and DESC (descending).
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) The name of the filter key.
* `values` - (Required, Set) One or more filter values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `event_buses` - event set information.
  * `add_time` - create time.
  * `connection_briefs` - Connector basic information, note: this field may return null, indicating that no valid value can be obtained.
    * `status` - Connector status, note: this field may return null, indicating that no valid value can be obtained.
    * `type` - Connector type, note: this field may return null, indicating that no valid value can be obtained.
  * `description` - Event set description, unlimited character type, description within 200 characters.
  * `event_bus_id` - event bus Id.
  * `event_bus_name` - Event set name, which can only contain letters, numbers, underscores, hyphens, starts with a letter and ends with a number or letter, 2~60 characters.
  * `mod_time` - update time.
  * `pay_mode` - Billing mode, note: this field may return null, indicating that no valid value can be obtained.
  * `target_briefs` - Target brief information, note: this field may return null, indicating that no valid value can be obtained.
    * `target_id` - Target ID.
    * `type` - Target type.
  * `type` - event bus type.


