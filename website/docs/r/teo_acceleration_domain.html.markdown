---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_acceleration_domain"
sidebar_current: "docs-tencentcloud-resource-teo_acceleration_domain"
description: |-
  Provides a resource to create a TEO acceleration domain
---

# tencentcloud_teo_acceleration_domain

Provides a resource to create a TEO acceleration domain

~> **NOTE:** Before modifying resource content, you need to ensure that the `status` is `online`.

## Example Usage

```hcl
resource "tencentcloud_teo_acceleration_domain" "example" {
  zone_id     = "zone-39quuimqg8r6"
  domain_name = "www.demo.com"

  origin_info {
    origin      = "150.109.8.1"
    origin_type = "IP_DOMAIN"
  }

  status            = "online"
  origin_protocol   = "FOLLOW"
  http_origin_port  = 80
  https_origin_port = 443
  ipv6_status       = "follow"
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
* `status` - (Optional, String) Accelerated domain name status, the values are: `online`: enabled; `offline`: disabled. Default is `online`.

The `origin_info` object supports the following:

* `origin_type` - (Required, String) The origin type. Values: `IP_DOMAIN`: IPv4/IPv6 address or domain name; `COS`: COS bucket address; `ORIGIN_GROUP`: Origin group; `AWS_S3`: AWS S3 bucket address; `SPACE`: EdgeOne Shield Space.
* `origin` - (Required, String) The origin address. Enter the origin group ID if `OriginType=ORIGIN_GROUP`.
* `backup_origin` - (Optional, String) ID of the secondary origin group (valid when `OriginType=ORIGIN_GROUP`). If it is not specified, it indicates that secondary origins are not used.
* `host_header` - (Optional, String) Customize the back-to-origin HOST header. This parameter is only valid when OriginType=IP_DOMAIN. If OriginType=COS or AWS_S3, the back-to-origin HOST header will be consistent with the origin server domain name. If OriginType=ORIGIN_GROUP, the back-to-origin HOST header follows the configuration in the origin server group. If no configuration is made, the default is the acceleration domain name. If OriginType=VOD or SPACE, there is no need to configure this header. It will take effect according to the corresponding back-to-origin domain name.
* `private_access` - (Optional, String) Whether to authenticate access to the private object storage origin (valid when `OriginType=COS/AWS_S3`). Values: `on`: Enable private authentication; `off`: Disable private authentication. If this field is not specified, the default value `off` is used.
* `private_parameters` - (Optional, List) The private authentication parameters. This field is valid when `PrivateAccess=on`.
* `vod_bucket_id` - (Optional, String) VOD bucket ID. This parameter is required when OriginType = VOD and VodOriginScope = bucket. Data source: the storage ID of the bucket in the Cloud VOD Professional Edition application.
* `vod_origin_scope` - (Optional, String) The scope of cloud on-demand back-to-source. This parameter is effective when OriginType = VOD. The possible values are: all: all files in the cloud on-demand application corresponding to the current origin station. The default value is all; bucket: files in a specified bucket under the cloud on-demand application corresponding to the current origin station. The bucket is specified by the parameter VodBucketId.

The `private_parameters` object of `origin_info` supports the following:

* `name` - (Required, String) The parameter name. Valid values: `AccessKeyId`: Access Key ID; `SecretAccessKey`: Secret Access Key; `SignatureVersion`: authentication version, v2 or v4; `Region`: bucket region.
* `value` - (Required, String) The parameter value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cname` - CNAME address.


## Import

TEO acceleration domain can be imported using the id, e.g.

```
terraform import tencentcloud_teo_acceleration_domain.example zone-39quuimqg8r6#www.demo.com
```

