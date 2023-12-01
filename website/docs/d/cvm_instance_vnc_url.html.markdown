---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cvm_instance_vnc_url"
sidebar_current: "docs-tencentcloud-datasource-cvm_instance_vnc_url"
description: |-
  Use this data source to query detailed information of cvm instance_vnc_url
---

# tencentcloud_cvm_instance_vnc_url

Use this data source to query detailed information of cvm instance_vnc_url

## Example Usage

```hcl
data "tencentcloud_cvm_instance_vnc_url" "instance_vnc_url" {
  instance_id = "ins-xxxxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID. To obtain the instance IDs, you can call `DescribeInstances` and look for `InstanceId` in the response.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_vnc_url` - Instance VNC URL.


