---
subcategory: "Waf"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_tls_versions"
sidebar_current: "docs-tencentcloud-datasource-waf_tls_versions"
description: |-
  Use this data source to query detailed information of waf tls_versions
---

# tencentcloud_waf_tls_versions

Use this data source to query detailed information of waf tls_versions

## Example Usage

```hcl
data "tencentcloud_waf_tls_versions" "example" {}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `tls` - TLS key value.
  * `version_id` - TLS version IDNote: This field may return null, indicating that a valid value cannot be obtained.
  * `version_name` - Tls version nameNote: This field may return null, indicating that a valid value cannot be obtained.


