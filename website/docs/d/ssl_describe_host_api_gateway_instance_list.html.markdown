---
subcategory: "SSL Certificates"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssl_describe_host_api_gateway_instance_list"
sidebar_current: "docs-tencentcloud-datasource-ssl_describe_host_api_gateway_instance_list"
description: |-
  Use this data source to query detailed information of ssl describe_host_api_gateway_instance_list
---

# tencentcloud_ssl_describe_host_api_gateway_instance_list

Use this data source to query detailed information of ssl describe_host_api_gateway_instance_list

## Example Usage

```hcl
data "tencentcloud_ssl_describe_host_api_gateway_instance_list" "describe_host_api_gateway_instance_list" {
  certificate_id = "9Bpk7XOu"
  resource_type  = "apiGateway"
}
```

## Argument Reference

The following arguments are supported:

* `certificate_id` - (Required, String) Certificate ID to be deployed.
* `resource_type` - (Required, String) Deploy resource type.
* `filters` - (Optional, List) List of filtering parameters; Filterkey: domainmatch.
* `is_cache` - (Optional, Int) Whether to query the cache, 1: Yes; 0: No, the default is the query cache, the cache is half an hour.
* `old_certificate_id` - (Optional, String) Deployed certificate ID.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `filter_key` - (Required, String) Filter parameter key.
* `filter_value` - (Required, String) Filter parameter value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_list` - Apigateway instance listNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `cert_id` - Certificate IDNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `domain` - domain name.
  * `protocol` - Use Agreement.
  * `service_id` - Instance ID.
  * `service_name` - Example name.


