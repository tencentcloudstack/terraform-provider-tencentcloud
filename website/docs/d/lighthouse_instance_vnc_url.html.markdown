---
subcategory: "TencentCloud Lighthouse(Lighthouse)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lighthouse_instance_vnc_url"
sidebar_current: "docs-tencentcloud-datasource-lighthouse_instance_vnc_url"
description: |-
  Use this data source to query detailed information of lighthouse instance_vnc_url
---

# tencentcloud_lighthouse_instance_vnc_url

Use this data source to query detailed information of lighthouse instance_vnc_url

## Example Usage

```hcl
data "tencentcloud_lighthouse_instance_vnc_url" "instance_vnc_url" {
  instance_id = "lhins-123456"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_vnc_url` - Instance VNC URL.


