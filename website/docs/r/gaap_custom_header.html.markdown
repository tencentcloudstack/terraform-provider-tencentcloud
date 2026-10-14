---
subcategory: "Global Application Acceleration(GAAP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_custom_header"
sidebar_current: "docs-tencentcloud-resource-gaap_custom_header"
description: |-
  Provides a resource to create a gaap custom_header
---

# tencentcloud_gaap_custom_header

Provides a resource to create a gaap custom_header

## Example Usage

```hcl
resource "tencentcloud_gaap_custom_header" "custom_header" {
  rule_id = "rule-xxxxxx"
  headers {
    header_name  = "HeaderName1"
    header_value = "HeaderValue1"
  }
  headers {
    header_name  = "HeaderName2"
    header_value = "HeaderValue2"
  }
}
```

## Argument Reference

The following arguments are supported:

* `rule_id` - (Required, String) Rule id.
* `headers` - (Optional, List) Headers.

The `headers` object supports the following:

* `header_name` - (Required, String) Header name.
* `header_value` - (Required, String) Header value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateCustomHeader](https://cloud.tencent.com/document/api/608/59242) - This API is used to create a custom header of the HTTP/HTTPS listener. When client requests reach the listener, they will be forwarded to the origin with this custom hearer.
* Read: [DescribeCustomHeader](https://cloud.tencent.com/document/api/608/59240) - This API is used to query the list of custom headers.
* Update: [CreateCustomHeader](https://cloud.tencent.com/document/api/608/59242) - This API is used to create a custom header of the HTTP/HTTPS listener. When client requests reach the listener, they will be forwarded to the origin with this custom hearer.
* Delete: [CreateCustomHeader](https://cloud.tencent.com/document/api/608/59242) - This API is used to create a custom header of the HTTP/HTTPS listener. When client requests reach the listener, they will be forwarded to the origin with this custom hearer.

## Import

gaap custom_header can be imported using the id, e.g.

```
terraform import tencentcloud_gaap_custom_header.custom_header ruleId
```

