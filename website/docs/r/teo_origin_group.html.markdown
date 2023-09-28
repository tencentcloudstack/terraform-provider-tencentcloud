---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_origin_group"
sidebar_current: "docs-tencentcloud-resource-teo_origin_group"
description: |-
  Provides a resource to create a teo origin_group
---

# tencentcloud_teo_origin_group

Provides a resource to create a teo origin_group

## Example Usage

### Self origin group

```hcl
resource "tencentcloud_teo_origin_group" "origin_group" {
  zone_id            = "zone-297z8rf93cfw"
  configuration_type = "weight"
  origin_group_name  = "test-group"
  origin_type        = "self"
  origin_records {
    area    = []
    port    = 8080
    private = false
    record  = "150.109.8.1"
    weight  = 100
  }
}
```

### Cos origin group

```hcl
resource "tencentcloud_teo_origin_group" "origin_group" {
  configuration_type = "weight"
  origin_group_name  = "test"
  origin_type        = "cos"
  zone_id            = "zone-2o3h21ed8bpu"

  origin_records {
    area    = []
    port    = 0
    private = true
    record  = "test-ruichaolin-1310708577.cos.ap-nanjing.myqcloud.com"
    weight  = 100
  }
}
```

## Argument Reference

The following arguments are supported:

* `configuration_type` - (Required, String) Type of the origin group, this field should be set when `OriginType` is self, otherwise leave it empty. Valid values: `area`: select an origin by using Geo info of the client IP and `Area` field in Records; `weight`: weighted select an origin by using `Weight` field in Records; `proto`: config by HTTP protocol.
* `origin_group_name` - (Required, String) OriginGroup Name.
* `origin_records` - (Required, List) Origin site records.
* `origin_type` - (Required, String) Type of the origin site. Valid values: `self`: self-build website; `cos`: tencent cos; `third_party`: third party cos.
* `zone_id` - (Required, String, ForceNew) Site ID.

The `origin_records` object supports the following:

* `port` - (Required, Int) Port of the origin site. Valid value range: 1-65535.
* `record` - (Required, String) Record value, which could be an IPv4/IPv6 address or a domain.
* `area` - (Optional, Set) Indicating origin sites area when `Type` field is `area`. An empty List indicate the default area. Valid value:- Asia, Americas, Europe, Africa or Oceania.
* `private_parameter` - (Optional, List) Parameters for private authentication. Only valid when `Private` is `true`.
* `private` - (Optional, Bool) Whether origin site is using private authentication. Only valid when `OriginType` is `third_party`.
* `weight` - (Optional, Int) Indicating origin sites weight when `Type` field is `weight`. Valid value range: 1-100. Sum of all weights should be 100.

The `private_parameter` object supports the following:

* `name` - (Required, String) Parameter Name. Valid values: `AccessKeyId`: Access Key ID; `SecretAccessKey`: Secret Access Key.
* `value` - (Required, String) Parameter value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `origin_group_id` - OriginGroup ID.
* `update_time` - Last modification date.


## Import

teo origin_group can be imported using the zone_id#originGroup_id, e.g.
````
terraform import tencentcloud_teo_origin_group.origin_group zone-297z8rf93cfw#origin-4f8a30b2-3720-11ed-b66b-525400dceb86
````

