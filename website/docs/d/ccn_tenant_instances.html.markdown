---
subcategory: "Cloud Connect Network(CCN)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ccn_tenant_instances"
sidebar_current: "docs-tencentcloud-datasource-ccn_tenant_instances"
description: |-
  Use this data source to query detailed information of vpc tenant_ccn
---

# tencentcloud_ccn_tenant_instances

Use this data source to query detailed information of vpc tenant_ccn

## Example Usage

```hcl
data "tencentcloud_ccn_tenant_instances" "tenant_ccn" {
  ccn_ids          = ["ccn-39lqkygf"]
  is_security_lock = ["true"]
}
```

## Argument Reference

The following arguments are supported:

* `ccn_ids` - (Optional, Set: [`String`]) filter by ccn ids, like: ['ccn-12345678'].
* `is_security_lock` - (Optional, Set: [`String`]) filter by locked, like ['true'].
* `result_output_file` - (Optional, String) Used to save results.
* `user_account_id` - (Optional, Set: [`String`]) filter by ccn ids, like: ['12345678'].


