---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_security_policy"
sidebar_current: "docs-tencentcloud-resource-gaap_security_policy"
description: |-
  Provides a resource to create a security policy of GAAP proxy.
---

# tencentcloud_gaap_security_policy

Provides a resource to create a security policy of GAAP proxy.

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
  action   = "DROP"
}
```

## Argument Reference

The following arguments are supported:

* `action` - (Required, ForceNew) Default policy, the available values includes `ACCEPT` and `DROP`.
* `proxy_id` - (Required, ForceNew) ID of the GAAP proxy.
* `enable` - (Optional) Indicates whether policy is enable, default is true.


## Import

GAAP security policy can be imported using the id, e.g.

```
  $ terraform import tencentcloud_gaap_security_policy.foo pl-xxxx
```

