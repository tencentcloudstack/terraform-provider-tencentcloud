---
subcategory: "DNSPod"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_subdomain_validate_status"
sidebar_current: "docs-tencentcloud-datasource-subdomain_validate_status"
description: |-
  Use this data source to query detailed information of dnspod subdomain_validate_status
---

# tencentcloud_subdomain_validate_status

Use this data source to query detailed information of dnspod subdomain_validate_status

## Example Usage

```hcl
data "tencentcloud_subdomain_validate_status" "subdomain_validate_status" {
  domain_zone = "www.iac-tf.cloud"
}
```

## Argument Reference

The following arguments are supported:

* `domain_zone` - (Required, String) Zone domain for which to view the verification status of TXT records.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `status` - Status. 0: not ready; 1: ready.


