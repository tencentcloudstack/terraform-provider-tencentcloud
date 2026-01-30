---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_l7_acc_rule_v2"
sidebar_current: "docs-tencentcloud-resource-teo_l7_acc_rule_v2"
description: |-
  Provides a resource to create a TEO l7 acc rule v2
---

# tencentcloud_teo_l7_acc_rule_v2

Provides a resource to create a TEO l7 acc rule v2

~> **NOTE:** Compared to tencentcloud_teo_l7_acc_rule, tencentcloud_teo_l7_acc_rule_v2 is simpler to use but is limited to managing a single rule and lacks the ability to maintain rule ordering. It is best suited for scenarios where you need to manage multiple rules independently and priority/sequencing is not a concern.

## Example Usage

```hcl
resource "tencentcloud_teo_l7_acc_rule_v2" "example" {
  zone_id     = "zone-3fkff38fyw8s"
  description = ["description"]
  rule_name   = "Web Acceleration"
  status      = "enable"
  branches {
    condition = "$${http.request.host} in ['www.example.com']"
    actions {
      name = "Cache"
      cache_parameters {
        custom_time {
          cache_time           = 2592000
          ignore_cache_control = "off"
          switch               = "on"
        }
      }
    }

    actions {
      name = "CacheKey"
      cache_key_parameters {
        full_url_cache = "on"
        ignore_case    = "off"
        query_string {
          switch = "off"
          values = []
        }
      }
    }

    actions {
      name = "ModifyRequestHeader"
      modify_request_header_parameters {
        header_actions {
          action = "set"
          name   = "EO-Client-OS"
          value  = "*"
        }

        header_actions {
          action = "add"
          name   = "O-Client-Browser"
          value  = "*"
        }

        header_actions {
          action = "del"
          name   = "Eo-Client-Device"
        }
      }
    }

    actions {
      name = "ContentCompression"
      content_compression_parameters {
        switch = "on"
      }
    }

    sub_rules {
      description = ["1-1"]
      branches {
        condition = "lower($${http.request.file_extension}) in ['php', 'jsp', 'asp', 'aspx']"
        actions {
          name = "Cache"
          cache_parameters {
            no_cache {
              switch = "on"
            }
          }
        }
      }
    }

    sub_rules {
      description = ["1-2"]
      branches {
        condition = "$${http.request.file_extension} in ['jpg', 'png', 'gif', 'bmp', 'svg', 'webp']"
        actions {
          name = "MaxAge"
          max_age_parameters {
            cache_time    = 3600
            follow_origin = "off"
          }
        }
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String, ForceNew) Zone id.
* `branches` - (Optional, List) Sub-Rule branch. this list currently supports filling in only one rule; multiple entries are invalid.
* `description` - (Optional, List: [`String`]) Rule annotation. multiple annotations can be added.
* `rule_name` - (Optional, String) Rule name. The name length limit is 255 characters.
* `status` - (Optional, String) Rule status. The possible values are: `enable`: enabled; `disable`: disabled.

The `access_url_redirect_parameters` object of `actions` supports the following:

* `host_name` - (Optional, List) Target hostname.
* `protocol` - (Optional, String) Target request protocol. valid values: http: target request protocol http; https: target request protocol https; follow: follow the request.
* `query_string` - (Optional, List) Carry query parameters.
* `status_code` - (Optional, Int) Status code. valid values: 301, 302, 303, 307, 308.
* `url_path` - (Optional, List) Target path.

The `actions` object of `branches` supports the following:

* `name` - (Required, String) Operation name. The name needs to correspond to the parameter structure, for example, if Name=Cache, CacheParameters is required.
- `Cache`: Node cache TTL;
- `CacheKey`: Custom Cache Key;
- `CachePrefresh`: Cache pre-refresh;
- `AccessURLRedirect`: Access URL redirection;
- `UpstreamURLRewrite`: Back-to-origin URL rewrite;
- `QUIC`: QUIC;
- `WebSocket`: WebSocket;
- `Authentication`: Token authentication;
- `MaxAge`: Browser cache TTL;
- `StatusCodeCache`: Status code cache TTL;
- `OfflineCache`: Offline cache;
- `SmartRouting`: Smart acceleration;
- `RangeOriginPull`: Segment back-to-origin;
- `UpstreamHTTP2`: HTTP2 back-to-origin;
- `HostHeader`: Host Header rewrite;
- `ForceRedirectHTTPS`: Access protocol forced HTTPS jump configuration;
- `OriginPullProtocol`: Back-to-origin HTTPS;
- `Compression`: Smart compression configuration;
- `HSTS`: HSTS;
- `ClientIPHeader`: Header information configuration for storing client request IP;
- `OCSPStapling`: OCSP stapling;
- `HTTP2`: HTTP2 Access;
- `PostMaxSize`: POST request upload file streaming maximum limit configuration;
- `ClientIPCountry`: Carry client IP region information when returning to the source;
- `UpstreamFollowRedirect`: Return to the source follow redirection parameter configuration;
- `UpstreamRequest`: Return to the source request parameters;
- `TLSConfig`: SSL/TLS security;
- `ModifyOrigin`: Modify the source station;
- `HTTPUpstreamTimeout`: Seven-layer return to the source timeout configuration;
- `HttpResponse`: HTTP response;
- `ErrorPage`: Custom error page;
- `ModifyResponseHeader`: Modify HTTP node response header;
- `ModifyRequestHeader`: Modify HTTP node request header;
- `ResponseSpeedLimit`: Single connection download speed limit.
- `SetContentIdentifierParameters`: Set content identifier.
* `access_url_redirect_parameters` - (Optional, List) The access url redirection configuration parameter. this parameter is required when name is accessurlredirect.
* `authentication_parameters` - (Optional, List) Token authentication configuration parameter. this parameter is required when name is authentication.
* `cache_key_parameters` - (Optional, List) Custom cache key configuration parameter. when name is cachekey, this parameter is required.
* `cache_parameters` - (Optional, List) Node cache ttl configuration parameter. when name is cache, this parameter is required.
* `cache_prefresh_parameters` - (Optional, List) The cache prefresh configuration parameter. this parameter is required when name is cacheprefresh.
* `client_ip_country_parameters` - (Optional, List) Configuration parameter for carrying the region information of the client ip during origin-pull. this parameter is required when the name is set to clientipcountry.
* `client_ip_header_parameters` - (Optional, List) Client ip header configuration for storing client request ip information. this parameter is required when name is clientipheader.
* `compression_parameters` - (Optional, List) Intelligent compression configuration. this parameter is required when name is set to compression.
* `content_compression_parameters` - (Optional, List) Content compression configuration parameters. This parameter is required when the `Name` parameter is set to `ContentCompression`. This parameter uses a whitelist function; please contact Tencent Cloud engineers if needed.
* `error_page_parameters` - (Optional, List) Custom error page configuration parameters. this parameter is required when name is errorpage.
* `force_redirect_https_parameters` - (Optional, List) Force https redirect configuration parameter. this parameter is required when the name is set to forceredirecthttps.
* `host_header_parameters` - (Optional, List) Host header rewrite configuration parameter. this parameter is required when name is set to hostheader.
* `hsts_parameters` - (Optional, List) HSTS configuration parameter. this parameter is required when name is hsts.
* `http2_parameters` - (Optional, List) HTTP2 access configuration parameter. this parameter is required when name is http2.
* `http_response_parameters` - (Optional, List) HTTP response configuration parameters. this parameter is required when name is httpresponse.
* `http_upstream_timeout_parameters` - (Optional, List) Configuration of layer 7 origin timeout. this parameter is required when name is httpupstreamtimeout.
* `max_age_parameters` - (Optional, List) Browser cache ttl configuration parameter. this parameter is required when name is maxage.
* `modify_origin_parameters` - (Optional, List) Configuration parameter for modifying the origin server. this parameter is required when the name is set to modifyorigin.
* `modify_request_header_parameters` - (Optional, List) Modify http node request header configuration parameters. this parameter is required when name is modifyrequestheader.
* `modify_response_header_parameters` - (Optional, List) Modify http node response header configuration parameters. this parameter is required when name is modifyresponseheader.
* `ocsp_stapling_parameters` - (Optional, List) OCSP stapling configuration parameter. this parameter is required when the name is set to ocspstapling.
* `offline_cache_parameters` - (Optional, List) Offline cache configuration parameter. this parameter is required when name is offlinecache.
* `origin_pull_protocol_parameters` - (Optional, List) Back-to-origin HTTPS configuration parameter. This parameter is required when the Name value is `OriginPullProtocol`.
* `post_max_size_parameters` - (Optional, List) Maximum size configuration for file streaming upload via a post request. this parameter is required when name is postmaxsize.
* `quic_parameters` - (Optional, List) The quic configuration parameter. this parameter is required when name is quic.
* `range_origin_pull_parameters` - (Optional, List) Shard source retrieval configuration parameter. this parameter is required when name is set to rangeoriginpull.
* `response_speed_limit_parameters` - (Optional, List) Single connection download speed limit configuration parameter. this parameter is required when name is responsespeedlimit.
* `set_content_identifier_parameters` - (Optional, List) Content identification configuration parameter. this parameter is required when name is httpresponse.
* `smart_routing_parameters` - (Optional, List) Smart acceleration configuration parameter. this parameter is required when name is smartrouting.
* `status_code_cache_parameters` - (Optional, List) Status code cache ttl configuration parameter. this parameter is required when name is statuscodecache.
* `tls_config_parameters` - (Optional, List) SSL/TLS security configuration parameter. this parameter is required when the name is set to tlsconfig.
* `upstream_follow_redirect_parameters` - (Optional, List) Configuration parameter for following redirects during origin-pull. this parameter is required when the name is set to upstreamfollowredirect.
* `upstream_http2_parameters` - (Optional, List) HTTP2 origin-pull configuration parameter. this parameter is required when name is set to upstreamhttp2.
* `upstream_request_parameters` - (Optional, List) Configuration parameter for origin-pull request. this parameter is required when the name is set to upstreamrequest.
* `upstream_url_rewrite_parameters` - (Optional, List) The origin-pull url rewrite configuration parameter. this parameter is required when name is upstreamurlrewrite.
* `web_socket_parameters` - (Optional, List) The websocket configuration parameter. this parameter is required when name is websocket.

The `authentication_parameters` object of `actions` supports the following:

* `auth_param` - (Optional, String) Authentication parameters name. the node will validate the value corresponding to this parameter name. consists of 1-100 uppercase and lowercase letters, numbers, or underscores.note: this field is required when authtype is either typea or typed.
* `auth_type` - (Optional, String) Authentication type. valid values:
- `TypeA`: authentication method a type, for specific meaning please refer to authentication method a. https://www.tencentcloud.com/document/product/1145/62475;
- `TypeB`: authentication method b type, for specific meaning please refer to authentication method b. https://www.tencentcloud.com/document/product/1145/62476;
- `TypeC`: authentication method c type, for specific meaning please refer to authentication method c. https://www.tencentcloud.com/document/product/1145/62477;
- `TypeD`: authentication method d type, for specific meaning please refer to authentication method d. https://www.tencentcloud.com/document/product/1145/62478;
- `TypeVOD`: authentication method v type, for specific meaning please refer to authentication method v. https://www.tencentcloud.com/document/product/1145/62479.
* `backup_secret_key` - (Optional, String) The backup authentication key consists of 6-40 uppercase and lowercase english letters or digits, and cannot contain " and $.
* `secret_key` - (Optional, String) The primary authentication key consists of 6-40 uppercase and lowercase english letters or digits, and cannot contain " and $.
* `time_format` - (Optional, String) Authentication time format. values: dec: decimal; hex: hexadecimal.
* `time_param` - (Optional, String) Authentication timestamp. it cannot be the same as the value of the authparam field.note: this field is required when authtype is typed.
* `timeout` - (Optional, Int) Validity period of the authentication url, in seconds, value range: 1-630720000. used to determine if the client access request has expired: If the current time exceeds "timestamp + validity period", it is an expired request, and a 403 is returned directly. If the current time does not exceed "timestamp + validity period", the request is not expired, and the md5 string is further validated. note: when authtype is one of typea, typeb, typec, or typed, this field is required.

The `branches` object of `sub_rules` supports the following:

* `actions` - (Optional, List) Sub-Rule branch. this list currently supports filling in only one rule; multiple entries are invalid.
* `condition` - (Optional, String) Match condition. https://www.tencentcloud.com/document/product/1145/54759.

The `branches` object of `sub_rules` supports the following:

* `actions` - (Optional, List) Sub-Rule branch. this list currently supports filling in only one rule; multiple entries are invalid.
* `condition` - (Optional, String) Match condition. https://www.tencentcloud.com/document/product/1145/54759.
* `sub_rules` - (Optional, List) List of sub-rules. multiple rules exist in this list and are executed sequentially from top to bottom. note: subrules and actions cannot both be empty. currently, only one layer of subrules is supported.

The `branches` object supports the following:

* `actions` - (Optional, List) Sub-Rule branch. this list currently supports filling in only one rule; multiple entries are invalid.
* `condition` - (Optional, String) Match condition. https://www.tencentcloud.com/document/product/1145/54759.
* `sub_rules` - (Optional, List) List of sub-rules. multiple rules exist in this list and are executed sequentially from top to bottom. note: subrules and actions cannot both be empty. currently, only one layer of subrules is supported.

The `cache_key_parameters` object of `actions` supports the following:

* `cookie` - (Optional, List) Cookie configuration parameters. at least one of the following configurations must be set: fullurlcache, ignorecase, header, scheme, cookie.
* `full_url_cache` - (Optional, String) Switch for retaining the complete query string. values: on: enable; off: disable.
* `header` - (Optional, List) HTTP request header configuration parameters. at least one of the following configurations must be set: fullurlcache, ignorecase, header, scheme, cookie.
* `ignore_case` - (Optional, String) Switch for ignoring case. values: enable; off: disable.note: at least one of fullurlcache, ignorecase, header, scheme, or cookie must be configured.
* `query_string` - (Optional, List) Configuration parameter for retaining the query string. this field and fullurlcache must be set simultaneously, but cannot both be on.
* `scheme` - (Optional, String) Request protocol switch. valid values: on: enable; off: disable.

The `cache_parameters` object of `actions` supports the following:

* `custom_time` - (Optional, List) Custom cache time. if not specified, this configuration is not set. only one of followorigin, nocache, or customtime can have switch set to on.
* `follow_origin` - (Optional, List) Cache follows origin server. if not specified, this configuration is not set. only one of followorigin, nocache, or customtime can have switch set to on.
* `no_cache` - (Optional, List) No cache. if not specified, this configuration is not set. only one of followorigin, nocache, or customtime can have switch set to on.

The `cache_prefresh_parameters` object of `actions` supports the following:

* `cache_time_percent` - (Optional, Int) Prefresh interval set as a percentage of the node cache time. value range: 1-99. note: this field is required when switch is on; when switch is off, this field is not required and will not take effect if filled.
* `switch` - (Optional, String) Whether to enable cache prefresh. values: enable; off: disable.

The `client_ip_country_parameters` object of `actions` supports the following:

* `header_name` - (Optional, String) Name of the request header that contains the client ip region. it is valid when switch=on. the default value eo-client-ipcountry is used when it is not specified.
* `switch` - (Optional, String) Whether to enable configuration. values: on: enable; off: disable.

The `client_ip_header_parameters` object of `actions` supports the following:

* `header_name` - (Optional, String) Name of the request header containing the client ip address for origin-pull. when switch is on, this parameter is required. x-forwarded-for is not allowed for this parameter.
* `switch` - (Optional, String) Whether to enable configuration. values: on: enable; off: disable.

The `compression_parameters` object of `actions` supports the following:

* `algorithms` - (Optional, List) Supported compression algorithm list. this field is required when switch is on; otherwise, it is not effective. valid values: brotli: brotli algorithm; gzip: gzip algorithm.
* `switch` - (Optional, String) Whether to enable smart compression. values: on: enable; off: disable.

The `content_compression_parameters` object of `actions` supports the following:

* `switch` - (Required, String) Content compression configuration switch, possible values are: on: enabled; off: disabled. When the Switch is set to `on`, both Brotli and gzip compression algorithms will be supported.

The `cookie` object of `cache_key_parameters` supports the following:

* `action` - (Optional, String) Cache action. values: full: retain all; ignore: ignore all; includeCustom: retain partial parameters; excludeCustom: ignore partial parameters. note: when switch is on, this field is required. when switch is off, this field is not required and will not take effect if filled.
* `switch` - (Optional, String) Whether to enable feature. values: on: enable; off: disable.
* `values` - (Optional, List) Custom cache key cookie name list.

The `cookie` object of `upstream_request_parameters` supports the following:

* `action` - (Optional, String) Origin-Pull request parameter cookie mode. this parameter is required when switch is on. valid values are: full: retain all; ignore: ignore all; includeCustom: retain partial parameters; excludeCustom: ignore partial parameters.
* `switch` - (Optional, String) Whether to enable the origin-pull request parameter cookie. valid values: on: enable; off: disable.
* `values` - (Optional, List) Specifies parameter values. this parameter takes effect only when the query string mode action is includecustom or excludecustom, and is used to specify the parameters to be reserved or ignored. up to 10 parameters are supported.

The `custom_time` object of `cache_parameters` supports the following:

* `cache_time` - (Optional, Int) Custom cache time value, unit: seconds. value range: 0-315360000.
* `ignore_cache_control` - (Optional, String) Ignore origin server cachecontrol switch. values: `on`: Enable; `off`: Disable.
* `switch` - (Optional, String) Custom cache time switch. values: `on`: Enable; `off`: Disable.

The `error_page_parameters` object of `actions` supports the following:

* `error_page_params` - (Optional, List) Custom error page configuration list.

The `error_page_params` object of `error_page_parameters` supports the following:

* `redirect_url` - (Required, String) Redirect url. requires a full redirect path, such as https://www.test.com/error.html.
* `status_code` - (Required, Int) Status code. supported values are 400, 403, 404, 405, 414, 416, 451, 500, 501, 502, 503, 504.

The `follow_origin` object of `cache_parameters` supports the following:

* `switch` - (Required, String) Whether to enable the configuration of following the origin server. Valid values: `on`: Enable; `off`: Disable.
* `default_cache_strategy` - (Optional, String) Whether to use the default caching policy when an origin server does not return the cache-control header. this field is required when defaultcache is set to on; otherwise, it is ineffective. when defaultcachetime is not 0, this field should be off. valid values: on: use the default caching policy. off: do not use the default caching policy.
* `default_cache_time` - (Optional, Int) The default cache time in seconds when an origin server does not return the cache-control header. the value ranges from 0 to 315360000. this field is required when defaultcache is set to on; otherwise, it is ineffective. when defaultcachestrategy is on, this field should be 0.
* `default_cache` - (Optional, String) Whether to cache when an origin server does not return the cache-control header. this field is required when switch is on; when switch is off, this field is not required and will be ineffective if filled. valid values: On: cache; Off: do not cache.

The `force_redirect_https_parameters` object of `actions` supports the following:

* `redirect_status_code` - (Optional, Int) Redirection status code. this field is required when switch is on; otherwise, it is not effective. valid values are: 301: 301 redirect; 302: 302 redirect.
* `switch` - (Optional, String) Whether to enable forced redirect configuration switch. values: on: enable; off: disable.

The `header_actions` object of `modify_request_header_parameters` supports the following:

* `action` - (Required, String) HTTP header setting methods. valid values are: set: sets a value for an existing header parameter; del: deletes a header parameter; add: adds a header parameter.
* `name` - (Required, String) HTTP header name.
* `value` - (Optional, String) HTTP header value. this parameter is required when the action is set to set or add; it is optional when the action is set to del.

The `header_actions` object of `modify_response_header_parameters` supports the following:

* `action` - (Required, String) HTTP header setting methods. valid values are: set: sets a value for an existing header parameter; del: deletes a header parameter; add: adds a header parameter.
* `name` - (Required, String) HTTP header name.
* `value` - (Optional, String) HTTP header value. this parameter is required when the action is set to set or add; it is optional when the action is set to del.

The `header` object of `cache_key_parameters` supports the following:

* `switch` - (Optional, String) Whether to enable feature. values: on: enable; off: disable.
* `values` - (Optional, List) Custom cache key http request header list. note: this field is required when switch is on; when switch is off, this field is not required and will not take effect if filled.

The `host_header_parameters` object of `actions` supports the following:

* `action` - (Optional, String) Action to be executed. values: followOrigin: follow origin server domain name; custom: custom.
* `server_name` - (Optional, String) Host header rewrite requires a complete domain name. note: this field is required when switch is on; when switch is off, this field is not required and any value will be ignored.

The `host_name` object of `access_url_redirect_parameters` supports the following:

* `action` - (Optional, String) Target hostname configuration, valid values are: follow: follow the request; custom: custom.
* `value` - (Optional, String) Custom value for target hostname, maximum length is 1024.

The `hsts_parameters` object of `actions` supports the following:

* `include_sub_domains` - (Optional, String) Whether to allow other subdomains to inherit the same hsts header. values: on: allows other subdomains to inherit the same hsts header; off: does not allow other subdomains to inherit the same hsts header. note: when switch is on, this field is required; when switch is off, this field is not required and will not take effect if filled.
* `preload` - (Optional, String) Whether to allow the browser to preload the hsts header. valid values: on: allows the browser to preload the hsts header; off: does not allow the browser to preload the hsts header. note: when switch is on, this field is required; when switch is off, this field is not required and will not take effect if filled.
* `switch` - (Optional, String) Whether to enable hsts. values: on: enable; off: disable.
* `timeout` - (Optional, Int) Cache hsts header time, unit: seconds. value range: 1-31536000. note: this field is required when switch is on; when switch is off, this field is not required and will not take effect if filled.

The `http2_parameters` object of `actions` supports the following:

* `switch` - (Optional, String) Whether to enable http2 access. values: on: enable; off: disable.

The `http_response_parameters` object of `actions` supports the following:

* `response_page` - (Optional, String) Response page id.
* `status_code` - (Optional, Int) Response status code. supports 2xx, 4xx, 5xx, excluding 499, 514, 101, 301, 302, 303, 509, 520-599.

The `http_upstream_timeout_parameters` object of `actions` supports the following:

* `response_timeout` - (Optional, Int) HTTP response timeout in seconds. value range: 5-600.

The `max_age_parameters` object of `actions` supports the following:

* `cache_time` - (Optional, Int) Custom cache time value, unit: seconds. value range: 0-315360000. note: when followorigin is off, it means not following the origin server and using cachetime to set the cache time; otherwise, this field will not take effect.
* `follow_origin` - (Optional, String) Specifies whether to follow the origin server cache-control configuration, with the following values: on: follow the origin server and ignore the field cachetime; off: do not follow the origin server and apply the field cachetime.

The `modify_origin_parameters` object of `actions` supports the following:

* `http_origin_port` - (Optional, Int) Ports for http origin-pull requests. value range: 1-65535. this parameter takes effect only when the origin-pull protocol originprotocol is http or follow.
* `https_origin_port` - (Optional, Int) Ports for https origin-pull requests. value range: 1-65535. this parameter takes effect only when the origin-pull protocol originprotocol is https or follow.
* `origin_protocol` - (Optional, String) Origin-Pull protocol configuration. this parameter is required when origintype is ipdomain, origingroup, or loadbalance. valid values are: Http: use http protocol; Https: use https protocol; Follow: follow the protocol.
* `origin_type` - (Optional, String) The origin type. values: IPDomain: ipv4, ipv6, or domain name type origin server; OriginGroup: origin server group type origin server; LoadBalance: cloud load balancer (clb), this feature is in beta test. to use it, please submit a ticket or contact smart customer service; COS: tencent cloud COS origin server; AWSS3: all object storage origin servers that support the aws s3 protocol.
* `origin` - (Optional, String) Origin server address, which varies according to the value of origintype: When origintype = ipdomain, fill in an ipv4 address, an ipv6 address, or a domain name; When origintype = cos, please fill in the access domain name of the cos bucket; When origintype = awss3, fill in the access domain name of the s3 bucket; When origintype = origingroup, fill in the origin server group id; When origintype = loadbalance, fill in the cloud load balancer instance id. this feature is currently only available to the allowlist.
* `private_access` - (Optional, String) Whether access to the private object storage origin server is allowed. this parameter is valid only when the origin server type origintype is COS or awss3. valid values: on: enable private authentication; off: disable private authentication. if not specified, the default value is off.
* `private_parameters` - (Optional, List) Private authentication parameter. this parameter is valid only when origintype = awss3 and privateaccess = on.

The `modify_request_header_parameters` object of `actions` supports the following:

* `header_actions` - (Optional, List) List of http header setting rules.

The `modify_response_header_parameters` object of `actions` supports the following:

* `header_actions` - (Optional, List) HTTP origin-pull header rules list.

The `no_cache` object of `cache_parameters` supports the following:

* `switch` - (Required, String) Whether to enable no-cache configuration. Valid values: `on`: Enable; `off`: Disable.

The `ocsp_stapling_parameters` object of `actions` supports the following:

* `switch` - (Optional, String) Whether to enable ocsp stapling configuration switch. values: on: enable; off: disable.

The `offline_cache_parameters` object of `actions` supports the following:

* `switch` - (Optional, String) Whether to enable offline caching. values: on: enable; Off: disable.

The `origin_pull_protocol_parameters` object of `actions` supports the following:

* `protocol` - (Optional, String) Back-to-origin protocol configuration. Possible values are: `http`: use HTTP protocol for back-to-origin; `https`: use HTTPS protocol for back-to-origin; `follow`: follow the protocol.

The `post_max_size_parameters` object of `actions` supports the following:

* `max_size` - (Optional, Int) Maximum size of the file uploaded for streaming via a post request, in bytes. value range: 1 * 2^20 bytes to 500 * 2^20 bytes.
* `switch` - (Optional, String) Whether to enable post request file upload limit, in bytes (default limit: 32 * 2^20 bytes). valid values: on: enable limit; off: disable limit.

The `private_parameters` object of `modify_origin_parameters` supports the following:

* `access_key_id` - (Required, String) Authentication parameter access key id.
* `secret_access_key` - (Required, String) Authentication parameter secret access key.
* `signature_version` - (Required, String) Authentication version. values: v2: v2 version; v4: v4 version.
* `region` - (Optional, String) Region of the bucket.

The `query_string` object of `access_url_redirect_parameters` supports the following:

* `action` - (Optional, String) Action to be executed. values: full: retain all; ignore: ignore all.

The `query_string` object of `cache_key_parameters` supports the following:

* `action` - (Optional, String) Actions to retain/ignore specified parameters in the query string. values: `includeCustom`: retain partial parameters. `excludeCustom`: ignore partial parameters.note: this field is required when switch is on. when switch is off, this field is not required and will not take effect if filled.
* `switch` - (Optional, String) Query string retain/ignore specified parameter switch. valid values are: on: enable; off: disable.
* `values` - (Optional, List) A list of parameter names to keep/ignore in the query string.

The `query_string` object of `upstream_request_parameters` supports the following:

* `action` - (Optional, String) Query string mode. this parameter is required when switch is on. values: full: retain all; ignore: ignore all; includeCustom: retain partial parameters; excludeCustom: ignore partial parameters.
* `switch` - (Optional, String) Whether to enable origin-pull request parameter query string. values: on: enable; off: disable.
* `values` - (Optional, List) Specifies parameter values. this parameter takes effect only when the query string mode action is includecustom or excludecustom, and is used to specify the parameters to be reserved or ignored. up to 10 parameters are supported.

The `quic_parameters` object of `actions` supports the following:

* `switch` - (Optional, String) Whether to enable quic. values: on: enable; off: disable.

The `range_origin_pull_parameters` object of `actions` supports the following:

* `switch` - (Optional, String) Whether to enable range gets. values are: on: enable; Off: disable.

The `response_speed_limit_parameters` object of `actions` supports the following:

* `max_speed` - (Required, String) Rate-Limiting value, in kb/s. enter a numerical value to specify the rate limit.
* `mode` - (Required, String) Download rate limit mode. valid values: LimitUponDownload: rate limit throughout the download process; LimitAfterSpecificBytesDownloaded: rate limit after downloading specific bytes at full speed; LimitAfterSpecificSecondsDownloaded: start speed limit after downloading at full speed for a specific duration.
* `start_at` - (Optional, String) Rate-Limiting start value, which can be the download size or specified duration, in kb or s. this parameter is required when mode is set to limitafterspecificbytesdownloaded or limitafterspecificsecondsdownloaded. enter a numerical value to specify the download size or duration.

The `set_content_identifier_parameters` object of `actions` supports the following:

* `content_identifier` - (Optional, String) Content identifier id.

The `smart_routing_parameters` object of `actions` supports the following:

* `switch` - (Optional, String) Whether to enable smart acceleration. values: on: enable; Off: disable.

The `status_code_cache_parameters` object of `actions` supports the following:

* `status_code_cache_params` - (Optional, List) Status code cache ttl.

The `status_code_cache_params` object of `status_code_cache_parameters` supports the following:

* `cache_time` - (Optional, Int) Cache time value in seconds. value range: 0-31536000.
* `status_code` - (Optional, Int) Status code. valid values: 400, 401, 403, 404, 405, 407, 414, 500, 501, 502, 503, 504, 509, 514.

The `sub_rules` object of `branches` supports the following:

* `branches` - (Optional, List) Sub-rule branch.
* `description` - (Optional, List) Rule comments.

The `tls_config_parameters` object of `actions` supports the following:

* `cipher_suite` - (Optional, String) Cipher suite. for detailed information, please refer to tls versions and cipher suites description, https://www.tencentcloud.com/document/product/1145/54154?has_map=1. valid values: loose-v2023: loose-v2023 cipher suite; general-v2023: general-v2023 cipher suite; strict-v2023: strict-v2023 cipher suite.
* `version` - (Optional, List) TLS version. at least one must be specified. if multiple versions are specified, they must be consecutive, e.g., enable tls1, 1.1, 1.2, and 1.3. it is not allowed to enable only 1 and 1.2 while disabling 1.1. valid values: tlsv1: tlsv1 version; `tlsv1.1`: tlsv1.1 version; `tlsv1.2`: tlsv1.2 version; `tlsv1.3`: tlsv1.3 version.

The `upstream_follow_redirect_parameters` object of `actions` supports the following:

* `max_times` - (Optional, Int) The maximum number of redirects. value range: 1-5. Note: this field is required when switch is on; when switch is off, this field is not required and will not take effect if filled.
* `switch` - (Optional, String) Whether to enable origin-pull to follow the redirection configuration. values: on: enable; off: disable.

The `upstream_http2_parameters` object of `actions` supports the following:

* `switch` - (Optional, String) Whether to enable http2 origin-pull. valid values: on: enable; off: disable.

The `upstream_request_parameters` object of `actions` supports the following:

* `cookie` - (Optional, List) Cookie configuration. optional. if not provided, it will not be configured.
* `query_string` - (Optional, List) Query string configuration. optional. if not provided, it will not be configured.

The `upstream_url_rewrite_parameters` object of `actions` supports the following:

* `action` - (Optional, String) Origin-Pull url rewrite action. valid values are: replace: replace the path prefix; addPrefix: add the path prefix; rmvPrefix: remove the path prefix.
* `regex` - (Optional, String) Origin URL Rewrite uses a regular expression for matching the complete path. It must conform to the Google RE2 specification and have a length range of 1 to 1024. This field is required when the Action is regexReplace; otherwise, it is optional.
* `type` - (Optional, String) Origin-Pull url rewriting type, only path is supported.
* `value` - (Optional, String) Origin-Pull url rewrite value, maximum length 1024, must start with /.note: when action is addprefix, it cannot end with /; when action is rmvprefix, * cannot be present.

The `url_path` object of `access_url_redirect_parameters` supports the following:

* `action` - (Optional, String) Action to be executed. values: follow: follow the request; custom: custom; regex: regular expression matching.
* `regex` - (Optional, String) Regular expression matching expression, length range is 1-1024. note: when action is regex, this field is required; when action is follow or custom, this field is not required and will not take effect if filled.
* `value` - (Optional, String) Redirect target url, length range is 1-1024.note: when action is regex or custom, this field is required; when action is follow, this field is not required and will not take effect if filled.

The `web_socket_parameters` object of `actions` supports the following:

* `switch` - (Optional, String) Whether to enable websocket connection timeout. values: on: use timeout as the websocket timeout;; off: the platform still supports websocket connections, using the system default timeout of 15 seconds.
* `timeout` - (Optional, Int) Timeout, unit: seconds. maximum timeout is 120 seconds.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `rule_id` - Rule ID. Unique identifier of the rule.
* `rule_priority` - Rule priority. only used as an output parameter.


## Import

TEO l7 acc rule v2 can be imported using the {zone_id}#{rule_id}, e.g.

````
terraform import tencentcloud_teo_l7_acc_rule_v2.example zone-3fkff38fyw8s#rule-3ft1xeuhlj1b
````

