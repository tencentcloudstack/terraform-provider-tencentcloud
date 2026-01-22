---
subcategory: "Intelligent Global Traffic Manager(IGTM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_igtm_detect_task_package_list"
sidebar_current: "docs-tencentcloud-datasource-igtm_detect_task_package_list"
description: |-
  Use this data source to query detailed information of IGTM detect task package list
---

# tencentcloud_igtm_detect_task_package_list

Use this data source to query detailed information of IGTM detect task package list

## Example Usage

### Query all igtm detect task package list

```hcl
data "tencentcloud_igtm_detect_task_package_list" "example" {}
```

### Query igtm detect task package list by filter

```hcl
data "tencentcloud_igtm_detect_task_package_list" "example" {
  filters {
    name  = "ResourceId"
    value = ["task-qqcoptejbwbf"]
    fuzzy = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Detect task filter conditions.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Filter field name, supported list as follows:
- ResourceId: detect task resource id.
- PeriodStart: minimum expiration time.
- PeriodEnd: maximum expiration time.
* `value` - (Required, Set) Filter field value.
* `fuzzy` - (Optional, Bool) Whether to enable fuzzy query, only supports filter field name as domain.
When fuzzy query is enabled, maximum Value length is 1, otherwise maximum Value length is 5. (Reserved field, not currently used).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `task_package_set` - Detect task package list.
  * `auto_renew_flag` - Whether auto-renew 0 no 1 yes.
  * `cost_item_list` - Billing item.
    * `cost_name` - Billing item name.
    * `cost_value` - Billing item value.
  * `create_time` - Package creation time.
  * `current_deadline` - Package expiration time.
  * `group` - Detect task type: 100 system setting; 200 billing; 300 management system; 110D monitoring migration free task; 120 disaster recovery switch task.
  * `is_expire` - Whether expired 0 no 1 yes.
  * `quota` - Quota.
  * `remark` - Remark.
  * `resource_id` - Resource ID.
  * `resource_type` - Resource type
TASK Detect task.
  * `status` - Status
ENABLED: Normal
ISOLATED: Isolated
DESTROYED: Destroyed
REFUNDED: Refunded.


