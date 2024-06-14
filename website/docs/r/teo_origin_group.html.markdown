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

~> **NOTE:** Please note that `tencentcloud_teo_origin_group` had to undergo incompatible changes in version v1.81.96.

## Example Usage

### Self origin group

```hcl
resource "tencentcloud_teo_origin_group" "basic" {
  name    = "keep-group-1"
  type    = "GENERAL"
  zone_id = "zone-197z8rf93cfw"

  records {
    record  = "tf-teo.xyz"
    type    = "IP_DOMAIN"
    weight  = 100
    private = true

    private_parameters {
      name  = "SecretAccessKey"
      value = "test"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `records` - (Required, Set) Origin site records.
* `type` - (Required, String) Type of the origin site. Valid values:
- `GENERAL`: Universal origin site group, only supports adding IP/domain name origin sites, which can be referenced by domain name service, rule engine, four-layer proxy, general load balancing, and HTTP-specific load balancing.
- `HTTP`: The HTTP-specific origin site group, supports adding IP/domain name and object storage origin site as the origin site, it cannot be referenced by the four-layer proxy, it can only be added to the acceleration domain name, rule engine-modify origin site, and HTTP-specific load balancing reference.
* `zone_id` - (Required, String, ForceNew) Site ID.
* `host_header` - (Optional, String) Back-to-origin Host Header, it only takes effect when type = HTTP is passed in. The rule engine modifies the Host Header configuration priority to be higher than the Host Header of the origin site group.
* `name` - (Optional, String) OriginGroup Name.

The `private_parameters` object of `records` supports the following:

* `name` - (Required, String) Private authentication parameter name, the values are:
  - `AccessKeyId`: Authentication parameter Access Key ID.
  - `SecretAccessKey`: Authentication parameter Secret Access Key.
  - `SignatureVersion`: Authentication version, v2 or v4.
  - `Region`: Bucket region.
* `value` - (Required, String) Private authentication parameter value.

The `records` object supports the following:

* `record` - (Required, String) Origin site record value, does not include port information, can be: IPv4, IPv6, domain name format.
* `private_parameters` - (Optional, List) Parameters for private authentication. Only valid when `Private` is `true`.
* `private` - (Optional, Bool) Whether to use private authentication, it takes effect when the origin site type RecordType=COS/AWS_S3, the values are:
  - `true`: Use private authentication.
  - `false`: Do not use private authentication.
* `record_id` - (Optional, String) Origin record ID.
* `type` - (Optional, String) Origin site type, the values are:
  - `IP_DOMAIN`: IPV4, IPV6, domain name type origin site.
  - `COS`: COS source.
  - `AWS_S3`: AWS S3 object storage origin site.
* `weight` - (Optional, Int) The weight of the origin site, the value is 0-100. If it is not filled in, it means that the weight will not be set and the system will schedule it freely. If it is filled in with 0, it means that the weight is 0 and the traffic will not be scheduled to this origin site.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Origin site group creation time.
* `origin_group_id` - OriginGroup ID.
* `references` - List of referenced instances of the origin site group.
  * `instance_id` - The instance ID of the reference type.
  * `instance_name` - Instance name of the application type.
  * `instance_type` - Reference service type, the values are:
  - `AccelerationDomain`: Acceleration domain name.
  - `RuleEngine`: Rule engine.
  - `Loadbalance`: Load balancing.
  - `ApplicationProxy`: Four-layer proxy.
* `update_time` - Origin site group update time.


## Import

teo origin_group can be imported using the zone_id#originGroup_id, e.g.
````
terraform import tencentcloud_teo_origin_group.origin_group zone-297z8rf93cfw#origin-4f8a30b2-3720-11ed-b66b-525400dceb86
````

