---
subcategory: "Serverless Cloud Function(SCF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_scf_custom_domain"
sidebar_current: "docs-tencentcloud-resource-scf_custom_domain"
description: |-
  Provides a resource to create a scf custom domain
---

# tencentcloud_scf_custom_domain

Provides a resource to create a scf custom domain

## Example Usage

```hcl
resource "tencentcloud_scf_custom_domain" "scf_custom_domain" {
  domain   = "xxxxxx"
  protocol = "HTTP"
  endpoints_config {
    namespace     = "default"
    function_name = "xxxxxx"
    qualifier     = "$LATEST"
    path_match    = "/aa/*"
  }
  waf_config {
    waf_open = "CLOSE"
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String) Domain names, pan-domain names are not supported.
* `endpoints_config` - (Required, List) Routing configuration.
* `protocol` - (Required, String) Protocol, value range: HTTP, HTTPS, HTTP&HTTPS.
* `cert_config` - (Optional, List) Certificate configuration information, required for HTTPS protocol.
* `waf_config` - (Optional, List) Web Application Firewall Configuration.

The `cert_config` object supports the following:

* `certificate_id` - (Optional, String) SSL Certificates ID.

The `endpoints_config` object supports the following:

* `function_name` - (Required, String) Function name.
* `namespace` - (Required, String) Function namespace.
* `path_match` - (Required, String) Path, value specification: /,/*,/xxx,/xxx/a,/xxx/*.
* `qualifier` - (Required, String) Function alias or version.
* `path_rewrite` - (Optional, List) Path rewriting policy.

The `path_rewrite` object of `endpoints_config` supports the following:

* `path` - (Required, String) Path that needs to be rerouted, value specification: /,/*,/xxx,/xxx/a,/xxx/*.
* `rewrite` - (Required, String) Replacement values: such as/, /$.
* `type` - (Required, String) Matching rules, value range: WildcardRules wildcard matching, ExactRules exact matching.

The `waf_config` object supports the following:

* `waf_instance_id` - (Optional, String) Web Application Firewall Instance ID.
* `waf_open` - (Optional, String) Whether the Web Application Firewall is turned on, value range:OPEN, CLOSE.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

scf scf_custom_domain can be imported using the id, e.g.

```
terraform import tencentcloud_scf_custom_domain.scf_custom_domain ${domain}
```

