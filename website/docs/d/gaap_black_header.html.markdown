---
subcategory: "Global Application Acceleration(GAAP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_black_header"
sidebar_current: "docs-tencentcloud-datasource-gaap_black_header"
description: |-
  Use this data source to query detailed information of gaap black header
---

# tencentcloud_gaap_black_header

Use this data source to query detailed information of gaap black header

## Example Usage

```hcl
data "tencentcloud_gaap_black_header" "black_header" {
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `black_headers` - Disabled custom header listNote: This field may return null, indicating that a valid value cannot be obtained.

## API Reference

This data source queries the following TencentCloud APIs:

* Read: [DescribeBlackHeader](https://cloud.tencent.com/document/api/608/59241) - This API is used to query names of blocked custom headers.

