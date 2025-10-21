---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_usable_unit_namespaces"
sidebar_current: "docs-tencentcloud-datasource-tsf_usable_unit_namespaces"
description: |-
  Use this data source to query detailed information of tsf usable_unit_namespaces
---

# tencentcloud_tsf_usable_unit_namespaces

Use this data source to query detailed information of tsf usable_unit_namespaces

## Example Usage

```hcl
data "tencentcloud_tsf_usable_unit_namespaces" "usable_unit_namespaces" {
  search_word = ""
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.
* `search_word` - (Optional, String) search by namespace id or namespace Name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - namespace object list.
  * `content` - namespace list.
    * `created_time` - Create time. Note: This field may return null, indicating that no valid value was found.
    * `gateway_instance_id` - Gateway instance id Note: This field may return null, indicating that no valid value was found.
    * `id` - Unit namespace ID. Note: This field may return null, indicating that no valid value was found.
    * `namespace_id` - namespace id.
    * `namespace_name` - namespace name.
    * `updated_time` - Update time. Note: This field may return null, indicating that no valid value was found.
  * `total_count` - total count.


