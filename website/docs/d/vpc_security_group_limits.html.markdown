---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_security_group_limits"
sidebar_current: "docs-tencentcloud-datasource-vpc_security_group_limits"
description: |-
  Use this data source to query detailed information of vpc security_group_limits
---

# tencentcloud_vpc_security_group_limits

Use this data source to query detailed information of vpc security_group_limits

## Example Usage

```hcl
data "tencentcloud_vpc_security_group_limits" "security_group_limits" {}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `security_group_limit_set` - sg limit set.
  * `instance_security_group_limit` - number of instances associated sg.
  * `referred_security_group_limit` - number of sg can be referred.
  * `security_group_extended_policy_limit` - number of sg extended policy.
  * `security_group_instance_limit` - number of sg associated instances.
  * `security_group_limit` - number of sg can be created.
  * `security_group_policy_limit` - number of sg polciy can be created.
  * `security_group_referred_cvm_and_eni_limit` - number of eni and cvm can be referred.
  * `security_group_referred_svc_limit` - number of svc can be referred.


