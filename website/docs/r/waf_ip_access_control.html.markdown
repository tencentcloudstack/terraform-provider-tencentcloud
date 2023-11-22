---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_ip_access_control"
sidebar_current: "docs-tencentcloud-resource-waf_ip_access_control"
description: |-
  Provides a resource to create a waf ip_access_control
---

# tencentcloud_waf_ip_access_control

Provides a resource to create a waf ip_access_control

## Example Usage

```hcl
resource "tencentcloud_waf_ip_access_control" "example" {
  instance_id = "waf_2kxtlbky00b3b4qz"
  domain      = "www.demo.com"
  edition     = "sparta-waf"
  items {
    ip       = "1.1.1.1"
    note     = "desc info."
    action   = 40
    valid_ts = "2019571199"
  }

  items {
    ip       = "2.2.2.2"
    note     = "desc info."
    action   = 42
    valid_ts = "2019571199"
  }

  items {
    ip       = "3.3.3.3"
    note     = "desc info."
    action   = 40
    valid_ts = "1680570420"
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String) Domain.
* `edition` - (Required, String) Waf edition. clb-waf means clb-waf, sparta-waf means saas-waf.
* `instance_id` - (Required, String) Waf instance Id.
* `items` - (Required, Set) Ip parameter list.

The `items` object supports the following:

* `action` - (Required, Int) Action value 40 is whitelist, 42 is blacklist.
* `ip` - (Required, String) IP address.
* `note` - (Required, String) Note info.
* `valid_ts` - (Required, Int) Effective date, with a second level timestamp value. For example, 1680570420 represents 2023-04-04 09:07:00; 2019571199 means permanently effective.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

waf ip_access_control can be imported using the id, e.g.

```
terraform import tencentcloud_waf_ip_access_control.example waf_2kxtlbky00b3b4qz#www.demo.com#sparta-waf
```

