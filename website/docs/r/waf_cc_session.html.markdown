---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_cc_session"
sidebar_current: "docs-tencentcloud-resource-waf_cc_session"
description: |-
  Provides a resource to create a waf cc_session
---

# tencentcloud_waf_cc_session

Provides a resource to create a waf cc_session

## Example Usage

```hcl
resource "tencentcloud_waf_cc_session" "example" {
  domain           = "www.demo.com"
  source           = "get"
  category         = "match"
  key_or_start_mat = "key_a=123"
  end_mat          = "&"
  start_offset     = "-1"
  end_offset       = "-1"
  edition          = "sparta-waf"
  session_name     = "terraformDemo"
}
```

## Argument Reference

The following arguments are supported:

* `category` - (Required, String) Session match pattern, Optional patterns are match, location.
* `domain` - (Required, String) Domain.
* `edition` - (Required, String) Waf edition. clb-waf means clb-waf, sparta-waf means saas-waf.
* `end_mat` - (Required, String) Session end identifier, when Category is match.
* `end_offset` - (Required, String) End offset position, when Category is location.
* `key_or_start_mat` - (Required, String) Session identifier.
* `session_name` - (Required, String) Session Name.
* `source` - (Required, String) Session matching position, Optional locations are get, post, header, cookie.
* `start_offset` - (Required, String) Starting offset position, when Category is location.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `session_id` - Session ID.


## Import

waf cc_session can be imported using the id, e.g.

```
terraform import tencentcloud_waf_cc_session.example www.demo.com#sparta-waf#2000000253
```

