---
subcategory: "Tencent Cloud Service Engine(TSE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tse_nacos_server_interfaces"
sidebar_current: "docs-tencentcloud-datasource-tse_nacos_server_interfaces"
description: |-
  Use this data source to query detailed information of tse nacos_server_interfaces
---

# tencentcloud_tse_nacos_server_interfaces

Use this data source to query detailed information of tse nacos_server_interfaces

## Example Usage

```hcl
data "tencentcloud_tse_nacos_server_interfaces" "nacos_server_interfaces" {
  instance_id = "ins-xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Optional, String) engine instance ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `content` - interface list.
  * `interface` - interface nameNote: This field may return null, indicating that a valid value is not available.


