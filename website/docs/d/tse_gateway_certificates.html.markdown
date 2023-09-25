---
subcategory: "Tencent Cloud Service Engine(TSE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tse_gateway_certificates"
sidebar_current: "docs-tencentcloud-datasource-tse_gateway_certificates"
description: |-
  Use this data source to query detailed information of tse gateway_certificates
---

# tencentcloud_tse_gateway_certificates

Use this data source to query detailed information of tse gateway_certificates

## Example Usage

```hcl
data "tencentcloud_tse_gateway_certificates" "gateway_certificates" {
  gateway_id = "gateway-ddbb709b"
  filters {
    key   = "BindDomain"
    value = "example.com"
  }
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, String) Gateway ID.
* `filters` - (Optional, List) Filter conditions, valid value: `BindDomain`, `Name`.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `key` - (Optional, String) Filter name.
* `value` - (Optional, String) Filter value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - Result.
  * `certificates_list` - Certificate list of gateway. Note: This field may return null, indicating that a valid value is not available.
    * `bind_domains` - Domains of the binding. Note: This field may return null, indicating that a valid value is not available.
    * `cert_id` - Certificate ID of ssl platform. Note: This field may return null, indicating that a valid value is not available.
    * `cert_source` - Source of certificate. Reference value:- native. Source: konga- ssl. Source: ssl platform. Note: This field may return null, indicating that a valid value is not available.
    * `create_time` - Upload time of certificate. Note: This field may return null, indicating that a valid value is not available.
    * `crt` - Pem format of certificate. Note: This field may return null, indicating that a valid value is not available.
    * `expire_time` - Expiration time of certificate. Note: This field may return null, indicating that a valid value is not available.
    * `id` - Certificate ID. Note: This field may return null, indicating that a valid value is not available.
    * `issue_time` - Issuance time of certificateNote: This field may return null, indicating that a valid value is not available.
    * `key` - Private key of certificate. Note: This field may return null, indicating that a valid value is not available.
    * `name` - Certificate name. Note: This field may return null, indicating that a valid value is not available.
    * `status` - Status of certificate. Reference value:- expired- active. Note: This field may return null, indicating that a valid value is not available.
  * `total` - Total count. Note: This field may return null, indicating that a valid value is not available.


