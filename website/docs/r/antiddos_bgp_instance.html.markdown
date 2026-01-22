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

~> **NOTE:** Currently, executing the `terraform destroy` command to delete this resource is not supported. If you need to destroy it, please contact Tencent Cloud AntiDDoS through a ticket.

## Example Usage

### Create standard bgp instance(POSTPAID)

```hcl
resource "tencentcloud_antiddos_bgp_instance" "example" {
  instance_charge_type = "POSTPAID_BY_MONTH"
  package_type         = "Standard"
  standard_package_config {
    region                 = "ap-guangzhou"
    protect_ip_count       = 1
    bandwidth              = 100
    elastic_bandwidth_flag = true
  }

  tag_info_list {
    tag_key   = "createBy"
    tag_value = "Terraform"
  }
}
```

### Create standard edition 2.0 bgp instance(PREPAID)

```hcl
resource "tencentcloud_antiddos_bgp_instance" "example" {
  instance_charge_type = "PREPAID"
  package_type         = "StandardPlus"
  instance_charge_prepaid {
    period     = 1
    renew_flag = "NOTIFY_AND_MANUAL_RENEW"
  }

  standard_plus_package_config {
    region                 = "ap-guangzhou"
    protect_count          = "TWO_TIMES"
    protect_ip_count       = 1
    bandwidth              = 100
    elastic_bandwidth_flag = true
  }

  tag_info_list {
    tag_key   = "createBy"
    tag_value = "Terraform"
  }
}
```

### Create enterprise bgp instance(POSTPAID)

```hcl
resource "tencentcloud_antiddos_bgp_instance" "example" {
  instance_charge_type = "POSTPAID_BY_MONTH"
  package_type         = "Enterprise"

  enterprise_package_config {
    region                  = "ap-guangzhou"
    protect_ip_count        = 10
    basic_protect_bandwidth = 300
    bandwidth               = 100
    elastic_bandwidth_flag  = false
  }

  tag_info_list {
    tag_key   = "createBy"
    tag_value = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_charge_type` - (Required, String, ForceNew) Payment Type: Payment Mode: PREPAID (Prepaid) / POSTPAID_BY_MONTH (Postpaid).
* `package_type` - (Required, String, ForceNew) High-defense package types: Enterprise, Standard, StandardPlus (Standard Edition 2.0).
* `enterprise_package_config` - (Optional, List, ForceNew) Enterprise package configuration.
* `instance_charge_prepaid` - (Optional, List, ForceNew) Prepaid configuration.
* `standard_package_config` - (Optional, List, ForceNew) Standard package configuration.
* `standard_plus_package_config` - (Optional, List, ForceNew) Standard Plus package configuration.
* `tag_info_list` - (Optional, List, ForceNew) Prepaid configuration.

The `enterprise_package_config` object supports the following:

* `bandwidth` - (Required, Int, ForceNew) Service bandwidth scale.
* `basic_protect_bandwidth` - (Required, Int, ForceNew) Guaranteed protection bandwidth.
* `protect_ip_count` - (Required, Int, ForceNew) Number of protected IPs.
* `region` - (Required, String, ForceNew) The region where the high-defense package was purchased.
* `elastic_bandwidth_flag` - (Optional, Bool, ForceNew) Whether to enable elastic service bandwidth. The default value is false.
* `elastic_protect_bandwidth` - (Optional, Int, ForceNew) Elastic bandwidth (Gbps), selectable elastic bandwidth [0, 400, 500, 600, 800, 1000], default is 0.

The `instance_charge_prepaid` object supports the following:

* `period` - (Optional, Int, ForceNew) Purchase period in months.
* `renew_flag` - (Optional, String, ForceNew) OTIFY_AND_MANUAL_RENEW: Notify the user of the expiration date and do not automatically renew. NOTIFY_AND_AUTO_RENEW: Notify the user of the expiration date and automatically renew. DISABLE_NOTIFY_AND_MANUAL_RENEW: Do not notify the user of the expiration date and do not automatically renew. The default is: Notify the user of the expiration date and do not automatically renew.

The `standard_package_config` object supports the following:

* `bandwidth` - (Required, Int, ForceNew) Protected service bandwidth 50Mbps.
* `protect_ip_count` - (Required, Int, ForceNew) Number of protected IPs.
* `region` - (Required, String, ForceNew) The region where the high-defense package was purchased.
* `elastic_bandwidth_flag` - (Optional, Bool, ForceNew) Whether to enable elastic service bandwidth. The default value is false.

The `standard_plus_package_config` object supports the following:

* `bandwidth` - (Required, Int, ForceNew) 50Mbps protected bandwidth.
* `protect_count` - (Required, String, ForceNew) Protection Count: TWO_TIMES: Two full-power protections; UNLIMITED: Infinite protections.
* `protect_ip_count` - (Required, Int, ForceNew) Number of protected IPs.
* `region` - (Required, String, ForceNew) The region where the high-defense package was purchased.
* `elastic_bandwidth_flag` - (Optional, Bool, ForceNew) Whether to enable elastic service bandwidth. The default value is false.

The `tag_info_list` object supports the following:

* `tag_key` - (Required, String, ForceNew) Tag key.
* `tag_value` - (Required, String, ForceNew) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `resource_id` - Bgp instance ID.


## Import

AntiDDoS bgp instance can be imported using the resourceId#packageRegion, e.g.

```
terraform import tencentcloud_antiddos_bgp_instance.example bgp-00000fyi#ap-guangzhou
```

