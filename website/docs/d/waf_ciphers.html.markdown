---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_ciphers"
sidebar_current: "docs-tencentcloud-datasource-waf_ciphers"
description: |-
  Use this data source to query detailed information of waf ciphers
---

# tencentcloud_waf_ciphers

Use this data source to query detailed information of waf ciphers

## Example Usage

```hcl
data "tencentcloud_waf_ciphers" "example" {}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `ciphers` - Encryption Suite InformationNote: This field may return null, indicating that a valid value cannot be obtained.
  * `cipher_id` - Encryption Suite IDNote: This field may return null, indicating that a valid value cannot be obtained.
  * `cipher_name` - Encryption Suite NameNote: This field may return null, indicating that a valid value cannot be obtained.
  * `version_id` - TLS version IDNote: This field may return null, indicating that a valid value cannot be obtained.

## API Reference

This data source queries the following TencentCloud APIs:

* Read: [DescribeCiphersDetail](https://cloud.tencent.com/document/api/627/96948) - 查询加密套件信息

