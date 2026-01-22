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

~> **NOTE:** Only `origin_type` is `IP_DOMAIN` can set `host_header`; And when `origin_type` is changed to `IP_DOMAIN`, `host_header` needs to be set to a legal value, such as a domain name string(like `domain_name`).

~> **NOTE:** If you use a third-party storage bucket configured for back-to-source, you need to ignore changes to `SecretAccessKey`.

~> **NOTE:** Before creating an accelerated domain, please verify the domain ownership of the site.

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

### Back-to-source configuration using a third-party storage bucket.

SecretAccessKey is sensitive data and can no longer be queried in plain text, so changes to SecretAccessKey need to be ignored.

```hcl
resource "tencentcloud_teo_acceleration_domain" "acceleration_domain" {
  domain_name       = "cos.demo.cn"
  http_origin_port  = 80
  https_origin_port = 443
  ipv6_status       = "follow"
  origin_protocol   = "FOLLOW"
  status            = "online"
  zone_id           = "zone-39quuimqg8r6"

  origin_info {
    backup_origin    = null
    origin           = "example.s3.ap-northeast.amazonaws.com"
    origin_type      = "AWS_S3"
    private_access   = "on"
    vod_bucket_id    = null
    vod_origin_scope = null

    private_parameters {
      name  = "AccessKeyId"
      value = "aaaaaaa"
    }
    private_parameters {
      name  = "SecretAccessKey"
      value = "bbbbbbb"
    }
    private_parameters {
      name  = "SignatureVersion"
      value = "v4"
    }
    private_parameters {
      name  = "Region"
      value = "us-east1"
    }
  }
  lifecycle {
    ignore_changes = [
      origin_info[0].private_parameters[1].value,
    ]
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
* `status` - (Optional, String) Accelerated domain name status, the values are: `online`: enabled; `offline`: disabled. Default is `online`.

The `origin_info` object supports the following:

* `origin_type` - (Required, String) Origin server type, with values: IP_DOMAIN: IPv4, IPv6, or domain name type origin server; COS: Tencent Cloud COS origin server; AWS_S3: AWS S3 origin server; ORIGIN_GROUP: origin server group type origin server; VOD: Video on Demand; SPACE: origin server uninstallation. Currently only available to the allowlist; LB: load balancing. Currently only available to the allowlist.
* `origin` - (Required, String) Origin server address, which varies according to the value of OriginType: When OriginType = IP_DOMAIN, fill in an IPv4 address, an IPv6 address, or a domain name; When OriginType = COS, fill in the access domain name of the COS bucket; When OriginType = AWS_S3, fill in the access domain name of the S3 bucket; When OriginType = ORIGIN_GROUP, fill in the origin server group ID; When OriginType = VOD, fill in the VOD application ID; When OriginType = LB, fill in the Cloud Load Balancer instance ID. This feature is currently only available to the allowlist; When OriginType = SPACE, fill in the origin server uninstallation space ID. This feature is currently only available to the allowlist.
* `backup_origin` - (Optional, String) The ID of the secondary origin group. This parameter is valid only when OriginType is ORIGIN_GROUP. This field indicates the old version capability, which cannot be configured or modified on the control panel after being called. Please submit a ticket if required.
* `host_header` - (Optional, String) Custom origin server HOST header. this parameter is valid only when OriginType=IP_DOMAIN.If the OriginType is another type of origin, this parameter does not need to be passed in, otherwise an error will be reported. If OriginType is COS or AWS_S3, the HOST header for origin-pull will remain consistent with the origin server domain name. If OriginType is ORIGIN_GROUP, the HOST header follows the ORIGIN site GROUP configuration. if not configured, it defaults to the acceleration domain name. If OriginType is VOD or SPACE, no configuration is required for this header, and the domain name takes effect based on the corresponding origin.
* `private_access` - (Optional, String) Whether access to the private Cloud Object Storage origin server is allowed. This parameter is valid only when OriginType is COS or AWS_S3. Valid values: on: Enable private authentication; off: Disable private authentication. If it is not specified, the default value is off.
* `private_parameters` - (Optional, List) Private authentication parameter. This parameter is valid only when `private_access` is on.
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

