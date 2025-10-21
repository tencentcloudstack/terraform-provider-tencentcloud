---
subcategory: "Serverless Cloud Function(SCF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_scf_triggers"
sidebar_current: "docs-tencentcloud-datasource-scf_triggers"
description: |-
  Use this data source to query detailed information of scf triggers
---

# tencentcloud_scf_triggers

Use this data source to query detailed information of scf triggers

## Example Usage

```hcl
data "tencentcloud_scf_triggers" "triggers" {
  function_name = "keep-1676351130"
  namespace     = "default"
  order_by      = "add_time"
  order         = "DESC"
}
```

## Argument Reference

The following arguments are supported:

* `function_name` - (Required, String) Function name.
* `filters` - (Optional, List) * Qualifier:Function version, alias.
* `namespace` - (Optional, String) Namespace. Default value: default.
* `order_by` - (Optional, String) Indicates by which field to sort the returned results. Valid values: add_time, mod_time. Default value: mod_time.
* `order` - (Optional, String) Indicates whether the returned results are sorted in ascending or descending order. Valid values: ASC, DESC. Default value: DESC.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Fields to be filtered. Up to 10 conditions allowed.Values of Name: VpcId, SubnetId, ClsTopicId, ClsLogsetId, Role, CfsId, CfsMountInsId, Eip. Values limit: 1.Name options: Status, Runtime, FunctionType, PublicNetStatus, AsyncRunEnable, TraceEnable. Values limit: 20.When Name is Runtime, CustomImage refers to the image type function.
* `values` - (Required, Set) Filter values of the field.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `triggers` - Trigger list.
  * `add_time` - Trigger creation time.
  * `available_status` - Whether the trigger is available.
  * `bind_status` - Trigger-Function binding status.
  * `custom_argument` - Custom parameterNote: this field may return null, indicating that no valid values can be obtained.
  * `enable` - Whether to enable.
  * `mod_time` - Trigger last modified time.
  * `qualifier` - Function version or alias.
  * `resource_id` - Minimum resource ID of trigger.
  * `trigger_attribute` - Trigger type. Two-way means that the trigger can be manipulated in both consoles, while one-way means that the trigger can be created only in the SCF Console.
  * `trigger_desc` - Detailed configuration of trigger.
  * `trigger_name` - Trigger name.
  * `type` - Trigger type.


