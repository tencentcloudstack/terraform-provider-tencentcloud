---
subcategory: "Serverless Cloud Function(SCF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_scf_function_aliases"
sidebar_current: "docs-tencentcloud-datasource-scf_function_aliases"
description: |-
  Use this data source to query detailed information of scf function_aliases
---

# tencentcloud_scf_function_aliases

Use this data source to query detailed information of scf function_aliases

## Example Usage

```hcl
data "tencentcloud_scf_function_aliases" "function_aliases" {
  function_name = "keep-1676351130"
  namespace     = "default"
}
```

## Argument Reference

The following arguments are supported:

* `function_name` - (Required, String) Function name.
* `function_version` - (Optional, String) If this parameter is provided, only aliases associated with this function version will be returned.
* `namespace` - (Optional, String) Function namespace.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `aliases` - Alias list.
  * `add_time` - Creation timeNote: this field may return null, indicating that no valid values can be obtained.
  * `description` - DescriptionNote: this field may return null, indicating that no valid values can be obtained.
  * `function_version` - Master version pointed to by the alias.
  * `mod_time` - Update timeNote: this field may return null, indicating that no valid values can be obtained.
  * `name` - Alias name.
  * `routing_config` - Routing information of aliasNote: this field may return null, indicating that no valid values can be obtained.
    * `addition_version_matchs` - Additional version with rule-based routing.
      * `expression` - Rule requirements for range match:It should be described in an open or closed range, i.e., `(a,b)` or `[a,b]`, where both a and b are integersRule requirements for exact match:Exact string match.
      * `key` - Matching rule key. When the API is called, pass in the `key` to route the request to the specified version based on the matching ruleHeader method:Enter invoke.headers.User for `key` and pass in `RoutingKey:{User:value}` when invoking a function through `invoke` for invocation based on rule matching.
      * `method` - Match method. Valid values:range: range matchexact: exact string match.
      * `version` - Function version name.
    * `additional_version_weights` - Additional version with random weight-based routing.
      * `version` - Function version name.
      * `weight` - Version weight.


