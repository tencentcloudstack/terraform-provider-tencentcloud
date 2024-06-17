---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_acceleration_domain"
sidebar_current: "docs-tencentcloud-resource-teo_acceleration_domain"
description: |-
  Provides a resource to create a teo acceleration_domain
---

# tencentcloud_teo_acceleration_domain

Provides a resource to create a teo acceleration_domain

## Example Usage

```hcl
resource "tencentcloud_teo_acceleration_domain" "acceleration_domain" {
  zone_id     = "zone-2o0i41pv2h8c"
  domain_name = "aaa.makn.cn"

  origin_info {
    origin      = "150.109.8.1"
    origin_type = "IP_DOMAIN"
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, String, ForceNew) Accelerated domain name.
* `origin_info` - (Required, List) Details of the origin.
* `zone_id` - (Required, String, ForceNew) ID of the site related with the accelerated domain name.
* `http_origin_port` - (Optional, Int) HTTP back-to-origin port, the value is 1-65535, effective when OriginProtocol=FOLLOW/HTTP, if not filled in, the default value is 80.
* `https_origin_port` - (Optional, Int) HTTPS back-to-origin port. The value range is 1-65535. It takes effect when OriginProtocol=FOLLOW/HTTPS. If it is not filled in, the default value is 443.
* `ipv6_status` - (Optional, String) IPv6 status, the value is: `follow`: follow the site IPv6 configuration; `on`: on; `off`: off. If not filled in, the default is: `follow`.
* `origin_protocol` - (Optional, String) Origin return protocol, possible values are: `FOLLOW`: protocol follow; `HTTP`: HTTP protocol back to source; `HTTPS`: HTTPS protocol back to source. If not filled in, the default is: `FOLLOW`.
* `status` - (Optional, String) Accelerated domain name status, the values are: `online`: enabled; `offline`: disabled.

The `origin_info` object supports the following:

* `origin_type` - (Required, String) The origin type. Values: `IP_DOMAIN`: IPv4/IPv6 address or domain name; `COS`: COS bucket address; `ORIGIN_GROUP`: Origin group; `AWS_S3`: AWS S3 bucket address; `SPACE`: EdgeOne Shield Space.
* `origin` - (Required, String) The origin address. Enter the origin group ID if `OriginType=ORIGIN_GROUP`.
* `backup_origin` - (Optional, String) ID of the secondary origin group (valid when `OriginType=ORIGIN_GROUP`). If it is not specified, it indicates that secondary origins are not used.
* `private_access` - (Optional, String) Whether to authenticate access to the private object storage origin (valid when `OriginType=COS/AWS_S3`). Values: `on`: Enable private authentication; `off`: Disable private authentication. If this field is not specified, the default value `off` is used.
* `private_parameters` - (Optional, List) The private authentication parameters. This field is valid when `PrivateAccess=on`.

The `private_parameters` object of `origin_info` supports the following:

* `name` - (Required, String) The parameter name. Valid values: `AccessKeyId`: Access Key ID; `SecretAccessKey`: Secret Access Key.
* `value` - (Required, String) The parameter value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cname` - CNAME address.


## Import

teo acceleration_domain can be imported using the id, e.g.

```
terraform import tencentcloud_teo_acceleration_domain.acceleration_domain acceleration_domain_id
```

