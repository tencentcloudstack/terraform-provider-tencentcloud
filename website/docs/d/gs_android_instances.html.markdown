---
subcategory: "GS"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gs_android_instances"
sidebar_current: "docs-tencentcloud-datasource-gs_android_instances"
description: |-
  Use this data source to query GS Android instances.
---

# tencentcloud_gs_android_instances

Use this data source to query GS Android instances.

## Example Usage

### Query all GS Android instances

```hcl
data "tencentcloud_gs_android_instances" "example" {}
```

### Query GS Android instances by filter

```hcl
data "tencentcloud_gs_android_instances" "example" {
  android_instance_ids = [
    "cai-1308726196-0352wk8np9s"
  ]
  android_instance_region = "ap-beijing"
  android_instance_zone   = "ap-beijing-1"
  label_selector {
    key      = "key"
    operator = "IN"
    values   = ["value"]
  }

  filters {
    name   = "Name"
    values = ["tf-example"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `android_instance_ids` - (Optional, List: [`String`]) List of instance IDs to query. Up to 100 per request.
* `android_instance_region` - (Optional, String) Instance region. Aggregated query across regions is not currently supported.
* `android_instance_zone` - (Optional, String) Instance availability zone.
* `filters` - (Optional, List) Field filters. Supported filter names: Name, UserId, HostSerialNumber, HostServerSerialNumber, AndroidInstanceModel.
* `label_selector` - (Optional, List) Instance label selector.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Filter field name.
* `values` - (Required, List) Filter field values.

The `label_selector` object supports the following:

* `key` - (Required, String) Label key.
* `operator` - (Required, String) Operator type. IN: label value must match one of Values; NOT_IN: must not match any; EXISTS: label key must exist; NOT_EXISTS: label key must not exist.
* `values` - (Optional, List) Label value list. Required for IN and NOT_IN operators.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `android_instance_list` - List of Android instances.
  * `android_instance_group_id` - Instance group ID.
  * `android_instance_id` - Instance ID.
  * `android_instance_image_id` - Instance image ID.
  * `android_instance_labels` - Instance label list.
    * `key` - Label key.
    * `value` - Label value.
  * `android_instance_model` - Android instance model. YS1: basic; GC0/GC1/GC2: performance.
  * `android_instance_region` - Instance region.
  * `android_instance_type` - Instance specification.
  * `android_instance_zone` - Instance availability zone.
  * `create_time` - Instance creation time.
  * `height` - Resolution height.
  * `host_serial_number` - Host serial number.
  * `host_server_serial_number` - Chassis serial number.
  * `name` - Instance name.
  * `private_ip` - Private IP address.
  * `service_status` - Service status. IDLE: not connected; ESTABLISHED: connected.
  * `state` - Instance state: INITIALIZING, NORMAL, PROCESSING.
  * `user_id` - User ID.
  * `width` - Resolution width.


