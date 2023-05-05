---
subcategory: "Serverless Cloud Function(SCF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_scf_function_alias"
sidebar_current: "docs-tencentcloud-resource-scf_function_alias"
description: |-
  Provides a resource to create a scf function_alias
---

# tencentcloud_scf_function_alias

Provides a resource to create a scf function_alias

## Example Usage

```hcl
// by weight
resource "tencentcloud_scf_function_alias" "function_alias" {
  description      = "weight test"
  function_name    = "keep-1676351130"
  function_version = "$LATEST"
  name             = "weight"
  namespace        = "default"

  routing_config {
    additional_version_weights {
      version = "2"
      weight  = 0.4
    }
  }
}

// by route
resource "tencentcloud_scf_function_alias" "function_alias" {
  description      = "matchs for test 12312312"
  function_name    = "keep-1676351130"
  function_version = "3"
  name             = "matchs"
  namespace        = "default"

  routing_config {
    additional_version_matches {
      expression = "testuser"
      key        = "invoke.headers.User"
      method     = "exact"
      version    = "2"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `function_name` - (Required, String) Function name.
* `function_version` - (Required, String) Master version pointed to by the alias.
* `name` - (Required, String) Alias name, which must be unique in the function, can contain 1 to 64 letters, digits, _, and -, and must begin with a letter.
* `description` - (Optional, String) Alias description information.
* `namespace` - (Optional, String) Function namespace.
* `routing_config` - (Optional, List) Request routing configuration of alias.

The `additional_version_matches` object supports the following:

* `expression` - (Required, String) Rule requirements for range match:It should be described in an open or closed range, i.e., (a,b) or [a,b], where both a and b are integersRule requirements for exact match:Exact string match.
* `key` - (Required, String) Matching rule key. When the API is called, pass in the key to route the request to the specified version based on the matching ruleHeader method:Enter invoke.headers.User for key and pass in RoutingKey:{User:value} when invoking a function through invoke for invocation based on rule matching.
* `method` - (Required, String) Match method. Valid values:range: Range matchexact: exact string match.
* `version` - (Required, String) Function version name.

The `additional_version_weights` object supports the following:

* `version` - (Required, String) Function version name.
* `weight` - (Required, Float64) Version weight.

The `routing_config` object supports the following:

* `additional_version_matches` - (Optional, List) Additional version with rule-based routing.
* `additional_version_weights` - (Optional, List) Additional version with random weight-based routing.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

scf function_alias can be imported using the id, e.g.

```
terraform import tencentcloud_scf_function_alias.function_alias namespace#functionName#name
```

