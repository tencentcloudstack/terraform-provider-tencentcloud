---
subcategory: "Anti-DDoS(DayuV2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_antiddos_ddos_geo_ip_block_config"
sidebar_current: "docs-tencentcloud-resource-antiddos_ddos_geo_ip_block_config"
description: |-
  Provides a resource to create a antiddos ddos_geo_ip_block_config
---

# tencentcloud_antiddos_ddos_geo_ip_block_config

Provides a resource to create a antiddos ddos_geo_ip_block_config

## Example Usage

```hcl
resource "tencentcloud_antiddos_ddos_geo_ip_block_config" "ddos_geo_ip_block_config" {
  instance_id = "bgp-xxxxxx"
  ddos_geo_ip_block_config {
    region_type = "customized"
    action      = "drop"
    area_list   = [100002]
  }
}
```

## Argument Reference

The following arguments are supported:

* `ddos_geo_ip_block_config` - (Required, List) DDoS region blocking configuration, configuration ID cannot be empty when filling in parameters.
* `instance_id` - (Required, String) InstanceId.

The `ddos_geo_ip_block_config` object supports the following:

* `action` - (Required, String) Blocking action, value [drop (intercept) trans (release)].
* `region_type` - (Required, String) Region type, value [oversea (overseas) China (domestic) customized (custom region)].
* `area_list` - (Optional, Set) When RegionType is customized, an AreaList must be filled in, with a maximum of 128 entries;.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

antiddos ddos_geo_ip_block_config can be imported using the id, e.g.

```
terraform import tencentcloud_antiddos_ddos_geo_ip_block_config.ddos_geo_ip_block_config ${instanceId}#${configId}
```

