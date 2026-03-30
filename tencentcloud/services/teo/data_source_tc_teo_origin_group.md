Use this data source to query detailed information of a TEO origin group

Example Usage

```hcl
data "tencentcloud_teo_origin_group" "foo" {
  zone_id = "zone-197z8rf93cfw"
}
```

Argument Reference

The following arguments are supported:

* `zone_id` - (Required, ForceNew) Site ID.

Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `origin_group_id` - OriginGroup ID.
* `name` - OriginGroup Name.
* `type` - Type of the origin site. Valid values:
    * `GENERAL`: Universal origin site group, only supports adding IP/domain name origin sites, which can be referenced by domain name service, rule engine, four-layer proxy, general load balancing, and HTTP-specific load balancing.
    * `HTTP`: The HTTP-specific origin site group, supports adding IP/domain name and object storage origin site as the origin site, it cannot be referenced by the four-layer proxy, it can only be added to the acceleration domain name, rule engine-modify origin site, and HTTP-specific load balancing reference.
* `records` - Origin site records.
    * `record` - Origin site record value, does not include port information, can be: IPv4, IPv6, domain name format.
    * `type` - Origin site type, the values are:
        * `IP_DOMAIN`: IPV4, IPV6, domain name type origin site.
        * `COS`: COS source.
        * `AWS_S3`: AWS S3 object storage origin site.
    * `record_id` - Origin record ID.
    * `weight` - The weight of the origin site, the value is 0-100. If it is not filled in, it means that the weight will not be set and the system will schedule it freely. If it is filled in with 0, it means that the weight is 0 and the traffic will not be scheduled to this origin site.
    * `private` - Whether to use private authentication, it takes effect when the origin site type RecordType=COS/AWS_S3, the values are:
        * `true`: Use private authentication.
        * `false`: Do not use private authentication.
    * `private_parameters` - Parameters for private authentication. Only valid when `Private` is `true`.
        * `name` - Private authentication parameter name, the values are:
            * `AccessKeyId`: Authentication parameter Access Key ID.
            * `SecretAccessKey`: Authentication parameter Secret Access Key.
            * `SignatureVersion`: Authentication version, v2 or v4.
            * `Region`: Bucket region.
        * `value` - Private authentication parameter value.
* `host_header` - Back-to-origin Host Header, it only takes effect when type = HTTP is passed in. The rule engine modifies the Host Header configuration priority to be higher than the Host Header of the origin site group.
* `references` - List of referenced instances of the origin site group.
    * `instance_type` - Reference service type, the values are:
        * `AccelerationDomain`: Acceleration domain name.
        * `RuleEngine`: Rule engine.
        * `Loadbalance`: Load balancing.
        * `ApplicationProxy`: Four-layer proxy.
    * `instance_id` - The instance ID of the reference type.
    * `instance_name` - Instance name of the application type.
* `create_time` - Origin site group creation time.
* `update_time` - Origin site group update time.
