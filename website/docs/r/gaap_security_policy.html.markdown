---
subcategory: "Global Application Acceleration(GAAP)"
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
  proxy_id = tencentcloud_gaap_proxy.foo.id
  action   = "DROP"
}
```

## Argument Reference

The following arguments are supported:

* `action` - (Required, String, ForceNew) Default policy. Valid value: `ACCEPT` and `DROP`.
* `proxy_id` - (Required, String, ForceNew) ID of the GAAP proxy.
* `enable` - (Optional, Bool) Indicates whether policy is enable, default value is `true`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateSecurityPolicy](https://cloud.tencent.com/document/api/608/36970)
* Create: [DescribeSecurityPolicyDetail](https://cloud.tencent.com/document/api/608/36959)
* Create: [OpenSecurityPolicy](https://cloud.tencent.com/document/api/608/36951)
* Read: [DescribeSecurityPolicyDetail](https://cloud.tencent.com/document/api/608/36959)
* Update: [CloseSecurityPolicy](https://cloud.tencent.com/document/api/608/36972)
* Update: [DescribeSecurityPolicyDetail](https://cloud.tencent.com/document/api/608/36959)
* Update: [OpenSecurityPolicy](https://cloud.tencent.com/document/api/608/36951)
* Delete: [CloseSecurityPolicy](https://cloud.tencent.com/document/api/608/36972)
* Delete: [DeleteSecurityPolicy](https://cloud.tencent.com/document/api/608/36968)
* Delete: [DescribeSecurityPolicyDetail](https://cloud.tencent.com/document/api/608/36959)

## Import

GAAP security policy can be imported using the id, e.g.

```
  $ terraform import tencentcloud_gaap_security_policy.foo pl-xxxx
```

