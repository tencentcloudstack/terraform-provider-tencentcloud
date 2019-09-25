---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_security_policies"
sidebar_current: "docs-tencentcloud-datasource-gaap_security_policies"
description: |-
  Use this data source to query security policies of GAAP proxy.
---

# tencentcloud_gaap_security_policies

Use this data source to query security policies of GAAP proxy.

## Example Usage

```hcl
resource "tencentcloud_gaap_proxy" "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

resource "tencentcloud_gaap_security_policy" "foo" {
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
  action   = "ACCEPT"
}

data "tencentcloud_gaap_security_policies" "foo" {
  id = "${tencentcloud_gaap_security_policy.foo.id}"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) ID of the security policy to be queried.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `action` - Default policy.
* `proxy_id` - ID of the GAAP proxy.
* `status` - Status of the security policy.


