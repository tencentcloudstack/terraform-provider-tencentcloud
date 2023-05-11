---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_delivery_config_by_group_id"
sidebar_current: "docs-tencentcloud-datasource-tsf_delivery_config_by_group_id"
description: |-
  Use this data source to query detailed information of tsf delivery_config_by_group_id
---

# tencentcloud_tsf_delivery_config_by_group_id

Use this data source to query detailed information of tsf delivery_config_by_group_id

## Example Usage

```hcl
data "tencentcloud_tsf_delivery_config_by_group_id" "delivery_config_by_group_id" {
  group_id = "group-yrjkln9v"
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, String) groupId.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - configuration item for deliver to a Kafka.
  * `config_id` - Config ID. Note: This field may return null, which means that no valid value was obtained.
  * `config_name` - Config Name. Note: This field may return null, which means that no valid value was obtained.


