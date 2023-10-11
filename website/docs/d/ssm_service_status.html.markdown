---
subcategory: "Secrets Manager(SSM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssm_service_status"
sidebar_current: "docs-tencentcloud-datasource-ssm_service_status"
description: |-
  Use this data source to query detailed information of ssm service_status
---

# tencentcloud_ssm_service_status

Use this data source to query detailed information of ssm service_status

## Example Usage

```hcl
data "tencentcloud_ssm_service_status" "example" {}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `access_key_escrow_enabled` - True means that the user can already use the key safe hosting function, false means that the user cannot use the key safe hosting function temporarily.
* `invalid_type` - Service unavailability type: 0-Not purchased, 1-Normal, 2-Service suspended due to arrears, 3-Resource release.
* `service_enabled` - True means the service has been activated, false means the service has not been activated yet.


