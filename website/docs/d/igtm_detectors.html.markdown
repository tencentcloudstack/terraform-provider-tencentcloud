---
subcategory: "Intelligent Global Traffic Manager(IGTM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_igtm_detectors"
sidebar_current: "docs-tencentcloud-datasource-igtm_detectors"
description: |-
  Use this data source to query detailed information of IGTM detectors
---

# tencentcloud_igtm_detectors

Use this data source to query detailed information of IGTM detectors

## Example Usage

```hcl
data "tencentcloud_igtm_detectors" "example" {}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `detector_group_set` - Detector group list.
  * `gid` - Line group ID GroupLineId.
  * `group_name` - Group name.
  * `group_type` - bgp, international, isp.
  * `internet_family` - ipv4, ipv6.
  * `package_set` - Supported package types.


