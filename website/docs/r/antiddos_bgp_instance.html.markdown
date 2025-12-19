---
subcategory: "Anti-DDoS(antiddos)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_antiddos_bgp_instance"
sidebar_current: "docs-tencentcloud-resource-antiddos_bgp_instance"
description: |-
  Provides a resource to create a AntiDDoS bgp instance
---

# tencentcloud_antiddos_bgp_instance

Provides a resource to create a AntiDDoS bgp instance

## Example Usage

```hcl
resource "tencentcloud_antiddos_bgp_instance" "example" {

}
```

## Argument Reference

The following arguments are supported:

* `instance_charge_type` - (Required, String) Payment Type: Payment Mode: PREPAID (Prepaid) / POSTPAID_BY_MONTH (Postpaid).
* `package_type` - (Required, String) High-defense package types: Enterprise, Standard, StandardPlus (Standard Edition 2.0).
* `enterprise_package_config` - (Optional, List) Enterprise package configuration.
* `instance_charge_prepaid` - (Optional, List) Prepaid configuration.
* `standard_package_config` - (Optional, List) Standard package configuration.
* `standard_plus_package_config` - (Optional, List) Standard Plus package configuration.
* `tag_info_list` - (Optional, List) Prepaid configuration.

The `enterprise_package_config` object supports the following:

* `bandwidth` - (Required, Int) Service bandwidth scale.
* `basic_protect_bandwidth` - (Required, Int) Guaranteed protection bandwidth.
* `protect_ip_count` - (Required, Int) Number of protected IPs.
* `region` - (Required, String) The region where the high-defense package was purchased.
* `elastic_bandwidth_flag` - (Optional, Bool) Whether to enable elastic service bandwidth. The default value is false.
* `elastic_protect_bandwidth` - (Optional, Int) Elastic bandwidth (Gbps), selectable elastic bandwidth [0, 400, 500, 600, 800, 1000], default is 0.

The `instance_charge_prepaid` object supports the following:

* `period` - (Optional, Int) Purchase period in months.
* `renew_flag` - (Optional, String) OTIFY_AND_MANUAL_RENEW: Notify the user of the expiration date and do not automatically renew. NOTIFY_AND_AUTO_RENEW: Notify the user of the expiration date and automatically renew. DISABLE_NOTIFY_AND_MANUAL_RENEW: Do not notify the user of the expiration date and do not automatically renew. The default is: Notify the user of the expiration date and do not automatically renew.

The `standard_package_config` object supports the following:

* `bandwidth` - (Required, Int) Protected service bandwidth 50Mbps.
* `protect_ip_count` - (Required, Int) Number of protected IPs.
* `region` - (Required, String) The region where the high-defense package was purchased.
* `elastic_bandwidth_flag` - (Optional, Bool) Whether to enable elastic service bandwidth. The default value is false.

The `standard_plus_package_config` object supports the following:

* `bandwidth` - (Required, Int) 50Mbps protected bandwidth.
* `protect_count` - (Required, String) Protection Count: TWO_TIMES: Two full-power protections; UNLIMITED: Infinite protections.
* `protect_ip_count` - (Required, Int) Number of protected IPs.
* `region` - (Required, String) The region where the high-defense package was purchased.
* `elastic_bandwidth_flag` - (Optional, Bool) Whether to enable elastic service bandwidth. The default value is false.

The `tag_info_list` object supports the following:

* `tag_key` - (Required, String) Tag key.
* `tag_value` - (Required, Int) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

AntiDDoS bgp instance can be imported using the id, e.g.

```
terraform import tencentcloud_antiddos_bgp_instance.example bgp-0000043i
```

