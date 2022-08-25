---
subcategory: "Teo"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_origin_group"
sidebar_current: "docs-tencentcloud-resource-teo_origin_group"
description: |-
  Provides a resource to create a teo originGroup
---

# tencentcloud_teo_origin_group

Provides a resource to create a teo originGroup

## Example Usage

```hcl
resource "tencentcloud_teo_origin_group" "originGroup" {
  origin_name = "test"
  type        = "weight"
  record {
    record  = "20160527-10003318.cos.ap-shanghai.myqcloud.com"
    area    = []
    weight  = 100
    port    = 0
    private = false

  }
  zone_id     = "zone-27mypfc1vr7d"
  origin_type = "cos"
}
```

## Argument Reference

The following arguments are supported:

* `origin_name` - (Required, String) OriginGroup Name.
* `record` - (Required, List) Origin website records.
* `type` - (Required, String) Type of the origin group, this field is required only when `OriginType` is `self`. Valid values:- area: select an origin by using Geo info of the client IP and `Area` field in Records.- weight: weighted select an origin by using `Weight` field in Records.
* `zone_id` - (Required, String) Site ID.
* `origin_type` - (Optional, String) Type of the origin website. Valid values:- self: self-build website.- cos: tencent cos.- third_party: third party cos.

The `private_parameter` object supports the following:

* `name` - (Required, String) Parameter Name. Valid values:- AccessKeyId: Access Key ID.- SecretAccessKey: Secret Access Key.
* `value` - (Required, String) Parameter value.

The `record` object supports the following:

* `area` - (Required, Set) Indicating origin website&#39;s area when `Type` field is `area`. An empty List indicate the default area.
* `port` - (Required, Int) Port of the origin website.
* `record` - (Required, String) Record Value.
* `weight` - (Required, Int) Indicating origin website&#39;s weight when `Type` field is `weight`. Valid value range: 1-100. Sum of all weights should be 100.
* `private_parameter` - (Optional, List) Parameters for private authentication. Only valid when `Private` is `true`.
* `private` - (Optional, Bool) Whether origin website is using private authentication. Only valid when `OriginType` is `third_party`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `update_time` - Last modification date.
* `zone_name` - Site Name.


## Import

teo originGroup can be imported using the id, e.g.
```
$ terraform import tencentcloud_teo_origin_group.originGroup zoneId#originId
```

