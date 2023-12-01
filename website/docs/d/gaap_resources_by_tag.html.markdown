---
subcategory: "Global Application Acceleration(GAAP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_resources_by_tag"
sidebar_current: "docs-tencentcloud-datasource-gaap_resources_by_tag"
description: |-
  Use this data source to query detailed information of gaap resources by tag
---

# tencentcloud_gaap_resources_by_tag

Use this data source to query detailed information of gaap resources by tag

## Example Usage

```hcl
data "tencentcloud_gaap_resources_by_tag" "resources_by_tag" {
  tag_key   = "tagKey"
  tag_value = "tagValue"
}
```

## Argument Reference

The following arguments are supported:

* `tag_key` - (Required, String) Tag key.
* `tag_value` - (Required, String) Tag value.
* `resource_type` - (Optional, String) Resource type, where:Proxy represents the proxy;ProxyGroup represents a proxy group;RealServer represents the Real Server.If this field is not specified, all resources under the label will be queried.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `resource_set` - List of resources corresponding to labels.
  * `resource_id` - Resource Id.
  * `resource_type` - Resource type, where:Proxy represents the proxy,ProxyGroup represents a proxy group,RealServer represents the real server.


