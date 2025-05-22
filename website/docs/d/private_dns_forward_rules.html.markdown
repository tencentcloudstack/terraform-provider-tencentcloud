---
subcategory: "PrivateDNS"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_private_dns_forward_rules"
sidebar_current: "docs-tencentcloud-datasource-private_dns_forward_rules"
description: |-
  Use this data source to query detailed information of Private Dns forward rules
---

# tencentcloud_private_dns_forward_rules

Use this data source to query detailed information of Private Dns forward rules

## Example Usage

### Query all private dns forward rules

```hcl
data "tencentcloud_private_dns_forward_rules" "example" {}
```

### Query all private dns forward rules by filters

```hcl
data "tencentcloud_private_dns_forward_rules" "example" {
  filters {
    name   = "RuleId"
    values = ["fid-2ece6ca305"]
  }

  filters {
    name   = "RuleName"
    values = ["tf-example"]
  }

  filters {
    name   = "RuleType"
    values = ["DOWN"]
  }

  filters {
    name   = "ZoneId"
    values = ["zone-04jlawty"]
  }

  filters {
    name   = "EndPointId"
    values = ["eid-e9d5880672"]
  }

  filters {
    name   = "EndPointName"
    values = ["tf-example"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter parameters.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Parameter name.
* `values` - (Required, Set) Array of parameter values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `forward_rule_set` - Private domain list.


