---
subcategory: "Waf"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_ports"
sidebar_current: "docs-tencentcloud-datasource-waf_ports"
description: |-
  Use this data source to query detailed information of waf ports
---

# tencentcloud_waf_ports

Use this data source to query detailed information of waf ports

## Example Usage

```hcl
data "tencentcloud_waf_ports" "example" {}
```

### Or

```hcl
data "tencentcloud_waf_ports" "example" {
  edition     = "clb-waf"
  instance_id = "waf_2kxtlbky00b2v1fn"
}
```

## Argument Reference

The following arguments are supported:

* `edition` - (Optional, String) Instance type, sparta-waf represents SAAS WAF, clb-waf represents CLB WAF.
* `instance_id` - (Optional, String) Instance unique ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `http_ports` - Http port list for instance.
* `https_ports` - Https port list for instance.


