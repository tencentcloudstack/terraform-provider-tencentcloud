---
subcategory: "API GateWay"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_upstreams"
sidebar_current: "docs-tencentcloud-datasource-api_gateway_upstreams"
description: |-
  Use this data source to query detailed information of apigateway upstream
---

# tencentcloud_api_gateway_upstreams

Use this data source to query detailed information of apigateway upstream

## Example Usage

```hcl
data "tencentcloud_api_gateway_upstreams" "example" {
  upstream_id = "upstream-4n5bfklc"
}
```

### Filtered Queries

```hcl
data "tencentcloud_api_gateway_upstreams" "example" {
  upstream_id = "upstream-4n5bfklc"

  filters {
    name   = "ServiceId"
    values = "service-hvg0uueg"
  }
}
```

## Argument Reference

The following arguments are supported:

* `upstream_id` - (Required, String) Backend channel ID.
* `filters` - (Optional, List) ServiceId and ApiId filtering queries.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Fields that need to be filtered.
* `values` - (Required, Set) The filtering value of the field.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - Query Results.
  * `api_id` - API Unique ID.
  * `api_name` - API nameNote: This field may return null, indicating that a valid value cannot be obtained.
  * `bind_time` - binding time.
  * `service_id` - Service Unique ID.
  * `service_name` - Service NameNote: This field may return null, indicating that a valid value cannot be obtained.


