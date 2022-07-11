---
subcategory: "Content Delivery Network(CDN)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cdn_domain"
sidebar_current: "docs-tencentcloud-resource-cdn_domain"
description: |-
  Provides a resource to create a CDN domain.
---

# tencentcloud_cdn_domain

Provides a resource to create a CDN domain.

## Example Usage

```hcl
resource "tencentcloud_cdn_domain" "foo" {
  domain         = "xxxx.com"
  service_type   = "web"
  area           = "mainland"
  full_url_cache = false

  origin {
    origin_type          = "ip"
    origin_list          = ["127.0.0.1"]
    origin_pull_protocol = "follow"
  }

  https_config {
    https_switch         = "off"
    http2_switch         = "off"
    ocsp_stapling_switch = "off"
    spdy_switch          = "off"
    verify_client        = "off"

    force_redirect {
      switch               = "on"
      redirect_type        = "http"
      redirect_status_code = 302
    }
  }

  tags = {
    hello = "world"
  }
}
```

Example Usage of cdn uses cache and request headers

```hcl
resource "tencentcloud_cdn_domain" "foo" {
  domain              = "xxxx.com"
  service_type        = "web"
  area                = "mainland"
  full_url_cache      = false
  range_origin_switch = "off"

  rule_cache {
    cache_time      = 10000
    no_cache_switch = "on"
    re_validate     = "on"
  }

  request_header {
    switch = "on"

    header_rules {
      header_mode  = "add"
      header_name  = "tf-header-name"
      header_value = "tf-header-value"
      rule_type    = "all"
      rule_paths   = ["*"]
    }
  }

  origin {
    origin_type          = "ip"
    origin_list          = ["127.0.0.1"]
    origin_pull_protocol = "follow"
  }

  https_config {
    https_switch         = "off"
    http2_switch         = "off"
    ocsp_stapling_switch = "off"
    spdy_switch          = "off"
    verify_client        = "off"

    force_redirect {
      switch               = "on"
      redirect_type        = "http"
      redirect_status_code = 302
    }
  }

  tags = {
    hello = "world"
  }
}
```

Example Usage of COS bucket url as origin

```hcl
resource "tencentcloud_cos_bucket" "bucket" {
  # Bucket format should be [custom name]-[appid].
  bucket = "demo-bucket-1251234567"
  acl    = "private"
}

# Create cdn domain
resource "tencentcloud_cdn_domain" "cdn" {
  domain         = "abc.com"
  service_type   = "web"
  area           = "mainland"
  full_url_cache = false

  origin {
    origin_type          = "cos"
    origin_list          = [tencentcloud_cos_bucket.bucket.cos_bucket_url]
    server_name          = tencentcloud_cos_bucket.bucket.cos_bucket_url
    origin_pull_protocol = "follow"
    cos_private_access   = "on"
  }

  https_config {
    https_switch         = "off"
    http2_switch         = "off"
    ocsp_stapling_switch = "off"
    spdy_switch          = "off"
    verify_client        = "off"
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) Name of the acceleration domain.
* `origin` - (Required, List) Origin server configuration. It's a list and consist of at most one item.
* `service_type` - (Required, String, ForceNew) Acceleration domain name service type. `web`: static acceleration, `download`: download acceleration, `media`: streaming media VOD acceleration.
* `area` - (Optional, String) Domain name acceleration region. `mainland`: acceleration inside mainland China, `overseas`: acceleration outside mainland China, `global`: global acceleration. Overseas acceleration service must be enabled to use overseas acceleration and global acceleration.
* `authentication` - (Optional, List) Specify timestamp hotlink protection configuration, NOTE: only one type can choose for the sub elements.
* `follow_redirect_switch` - (Optional, String) 301/302 redirect following switch, available values: `on`, `off` (default).
* `full_url_cache` - (Optional, Bool) Whether to enable full-path cache. Default value is `true`.
* `https_config` - (Optional, List) HTTPS acceleration configuration. It's a list and consist of at most one item.
* `ipv6_access_switch` - (Optional, String) ipv6 access configuration switch. Only available when area set to `mainland`. Valid values are `on` and `off`. Default value is `off`.
* `project_id` - (Optional, Int) The project CDN belongs to, default to 0.
* `range_origin_switch` - (Optional, String) Sharding back to source configuration switch. Valid values are `on` and `off`. Default value is `on`.
* `request_header` - (Optional, List) Request header configuration. It's a list and consist of at most one item.
* `rule_cache` - (Optional, List) Advanced path cache configuration.
* `tags` - (Optional, Map) Tags of cdn domain.

The `authentication` object supports the following:

* `switch` - (Optional, String) Authentication switching, available values: `on`, `off`.
* `type_a` - (Optional, List) Timestamp hotlink protection mode A configuration.
* `type_b` - (Optional, List) Timestamp hotlink protection mode B configuration. NOTE: according to upgrading of TencentCloud Platform, TypeB is unavailable for now.
* `type_c` - (Optional, List) Timestamp hotlink protection mode C configuration.
* `type_d` - (Optional, List) Timestamp hotlink protection mode D configuration.

The `client_certificate_config` object supports the following:

* `certificate_content` - (Required, String) Client Certificate PEM format, requires Base64 encoding.

The `force_redirect` object supports the following:

* `redirect_status_code` - (Optional, Int) Forced redirect status code. Valid values are `301` and `302`. When `switch` setting `off`, this property does not need to be set or set to `302`. Default value is `302`.
* `redirect_type` - (Optional, String) Forced redirect type. Valid values are `http` and `https`. `http` means a forced redirect from HTTPS to HTTP, `https` means a forced redirect from HTTP to HTTPS. When `switch` setting `off`, this property does not need to be set or set to `http`. Default value is `http`.
* `switch` - (Optional, String) Forced redirect configuration switch. Valid values are `on` and `off`. Default value is `off`.

The `header_rules` object supports the following:

* `header_mode` - (Required, String) Http header setting method. The following types are supported: `add`: add a head, if a head already exists, there will be a duplicate head, `del`: delete the head.
* `header_name` - (Required, String) Http header name.
* `header_value` - (Required, String) Http header value, optional when Mode is `del`, Required when Mode is `add`/`set`.
* `rule_paths` - (Required, List) Matching content under the corresponding type of CacheType: `all`: fill *, `file`: fill in the suffix name, such as jpg, txt, `directory`: fill in the path, such as /xxx/test, `path`: fill in the absolute path, such as /xxx/test.html.
* `rule_type` - (Required, String) Rule type. The following types are supported: `all`: all documents take effect, `file`: the specified file suffix takes effect, `directory`: the specified path takes effect, `path`: specify the absolute path to take effect.

The `https_config` object supports the following:

* `https_switch` - (Required, String) HTTPS configuration switch. Valid values are `on` and `off`.
* `client_certificate_config` - (Optional, List) Client certificate configuration information.
* `force_redirect` - (Optional, List) Configuration of forced HTTP or HTTPS redirects.
* `http2_switch` - (Optional, String) HTTP2 configuration switch. Valid values are `on` and `off`. and default value is `off`.
* `ocsp_stapling_switch` - (Optional, String) OCSP configuration switch. Valid values are `on` and `off`. and default value is `off`.
* `server_certificate_config` - (Optional, List) Server certificate configuration information.
* `spdy_switch` - (Optional, String) Spdy configuration switch. Valid values are `on` and `off`. and default value is `off`. This parameter is for white-list customer.
* `verify_client` - (Optional, String) Client certificate authentication feature. Valid values are `on` and `off`. and default value is `off`.

The `origin` object supports the following:

* `origin_list` - (Required, List) Master origin server list. Valid values can be ip or domain name. When modifying the origin server, you need to enter the corresponding `origin_type`.
* `origin_type` - (Required, String) Master origin server type. The following types are supported: `domain`: domain name type, `cos`: COS origin, `ip`: IP list used as origin server, `ipv6`: origin server list is a single IPv6 address, `ip_ipv6`: origin server list is multiple IPv4 addresses and an IPv6 address.
* `backup_origin_list` - (Optional, List) Backup origin server list. Valid values can be ip or domain name. When modifying the backup origin server, you need to enter the corresponding `backup_origin_type`.
* `backup_origin_type` - (Optional, String) Backup origin server type, which supports the following types: `domain`: domain name type, `ip`: IP list used as origin server.
* `backup_server_name` - (Optional, String) Host header used when accessing the backup origin server. If left empty, the ServerName of master origin server will be used by default.
* `cos_private_access` - (Optional, String) When OriginType is COS, you can specify if access to private buckets is allowed. Valid values are `on` and `off`. and default value is `off`.
* `origin_pull_protocol` - (Optional, String) Origin-pull protocol configuration. `http`: forced HTTP origin-pull, `follow`: protocol follow origin-pull, `https`: forced HTTPS origin-pull. This only supports origin server port 443 for origin-pull.
* `server_name` - (Optional, String) Host header used when accessing the master origin server. If left empty, the acceleration domain name will be used by default.

The `request_header` object supports the following:

* `header_rules` - (Optional, List) Custom request header configuration rules.
* `switch` - (Optional, String) Custom request header configuration switch. Valid values are `on` and `off`. and default value is `off`.

The `rule_cache` object supports the following:

* `cache_time` - (Required, Int) Cache expiration time setting, the unit is second, the maximum can be set to 365 days.
* `compare_max_age` - (Optional, String) Advanced cache expiration configuration. When it is turned on, it will compare the max-age value returned by the origin site with the cache expiration time set in CacheRules, and take the minimum value to cache at the node. Valid values are `on` and `off`. Default value is `off`.
* `follow_origin_switch` - (Optional, String) Follow the source station configuration switch. Valid values are `on` and `off`.
* `ignore_cache_control` - (Optional, String) Force caching. After opening, the no-store and no-cache resources returned by the origin site will also be cached in accordance with the CacheRules rules. Valid values are `on` and `off`. Default value is `off`.
* `ignore_set_cookie` - (Optional, String) Ignore the Set-Cookie header of the origin site. Valid values are `on` and `off`. Default value is `off`. This parameter is for white-list customer.
* `no_cache_switch` - (Optional, String) Cache configuration switch. Valid values are `on` and `off`.
* `re_validate` - (Optional, String) Always check back to origin. Valid values are `on` and `off`. Default value is `off`.
* `rule_paths` - (Optional, List) Matching content under the corresponding type of CacheType: `all`: fill *, `file`: fill in the suffix name, such as jpg, txt, `directory`: fill in the path, such as /xxx/test, `path`: fill in the absolute path, such as /xxx/test.html, `index`: fill /, `default`: Fill `no max-age`.
* `rule_type` - (Optional, String) Rule type. The following types are supported: `all`: all documents take effect, `file`: the specified file suffix takes effect, `directory`: the specified path takes effect, `path`: specify the absolute path to take effect, `index`: home page, `default`: effective when the source site has no max-age.
* `switch` - (Optional, String) Cache configuration switch. Valid values are `on` and `off`.

The `server_certificate_config` object supports the following:

* `certificate_content` - (Optional, String) Server certificate information. This is required when uploading an external certificate, which should contain the complete certificate chain.
* `certificate_id` - (Optional, String) Server certificate ID.
* `message` - (Optional, String) Certificate remarks.
* `private_key` - (Optional, String) Server key information. This is required when uploading an external certificate.

The `type_a` object supports the following:

* `expire_time` - (Required, Int) Signature expiration time in second. The maximum value is 630720000.
* `file_extensions` - (Required, List) File extension list settings determining if authentication should be performed. NOTE: If it contains an asterisk (*), this indicates all files.
* `filter_type` - (Required, String) Available values: `whitelist` - all types apart from `file_extensions` are authenticated, `blacklist`: - only the types in the `file_extensions` are authenticated.
* `secret_key` - (Required, String) The key for signature calculation. Only digits, upper and lower-case letters are allowed. Length limit: 6-32 characters.
* `sign_param` - (Required, String) Signature parameter name. Only upper and lower-case letters, digits, and underscores (_) are allowed. It cannot start with a digit. Length limit: 1-100 characters.
* `backup_secret_key` - (Optional, String) Used for calculate a signature. 6-32 characters. Only digits and letters are allowed.

The `type_b` object supports the following:

* `expire_time` - (Required, Int) Signature expiration time in second. The maximum value is 630720000.
* `file_extensions` - (Required, List) File extension list settings determining if authentication should be performed. NOTE: If it contains an asterisk (*), this indicates all files.
* `filter_type` - (Required, String) Available values: `whitelist` - all types apart from `file_extensions` are authenticated, `blacklist`: - only the types in the `file_extensions` are authenticated.
* `secret_key` - (Required, String) The key for signature calculation. Only digits, upper and lower-case letters are allowed. Length limit: 6-32 characters.
* `backup_secret_key` - (Optional, String) Used for calculate a signature. 6-32 characters. Only digits and letters are allowed.

The `type_c` object supports the following:

* `expire_time` - (Required, Int) Signature expiration time in second. The maximum value is 630720000.
* `file_extensions` - (Required, List) File extension list settings determining if authentication should be performed. NOTE: If it contains an asterisk (*), this indicates all files.
* `filter_type` - (Required, String) Available values: `whitelist` - all types apart from `file_extensions` are authenticated, `blacklist`: - only the types in the `file_extensions` are authenticated.
* `secret_key` - (Required, String) The key for signature calculation. Only digits, upper and lower-case letters are allowed. Length limit: 6-32 characters.
* `backup_secret_key` - (Optional, String) Used for calculate a signature. 6-32 characters. Only digits and letters are allowed.
* `time_format` - (Optional, String) Timestamp formation, available values: `dec`, `hex`.

The `type_d` object supports the following:

* `expire_time` - (Required, Int) Signature expiration time in second. The maximum value is 630720000.
* `file_extensions` - (Required, List) File extension list settings determining if authentication should be performed. NOTE: If it contains an asterisk (*), this indicates all files.
* `filter_type` - (Required, String) Available values: `whitelist` - all types apart from `file_extensions` are authenticated, `blacklist`: - only the types in the `file_extensions` are authenticated.
* `secret_key` - (Required, String) The key for signature calculation. Only digits, upper and lower-case letters are allowed. Length limit: 6-32 characters.
* `backup_secret_key` - (Optional, String) Used for calculate a signature. 6-32 characters. Only digits and letters are allowed.
* `time_format` - (Optional, String) Timestamp formation, available values: `dec`, `hex`.
* `time_param` - (Optional, String) Timestamp parameter name. Only upper and lower-case letters, digits, and underscores (_) are allowed. It cannot start with a digit. Length limit: 1-100 characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cname` - CNAME address of domain name.
* `create_time` - Creation time of domain name.
* `status` - Acceleration service status.


## Import

CDN domain can be imported using the id, e.g.

```
$ terraform import tencentcloud_cdn_domain.foo xxxx.com
```

