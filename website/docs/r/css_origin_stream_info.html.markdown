---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_origin_stream_info"
sidebar_current: "docs-tencentcloud-resource-css_origin_stream_info"
description: |-
  Provides a resource to create a CSS origin stream info.
---

# tencentcloud_css_origin_stream_info

Provides a resource to create a CSS origin stream info.

## Example Usage

```hcl
resource "tencentcloud_css_origin_stream_info" "example" {
  domain_name             = "www.demo.com"
  origin_stream_play_type = "rtmp"
  cdn_stream_play_type    = ["rtmp"]
  origin_stream_type      = 1
  origin_address_type     = 1
  origin_address          = ["1.1.1.1:8080"]
  origin_timeout          = 10000
  origin_retry_times      = 10
}
```

## Argument Reference

The following arguments are supported:

* `cdn_stream_play_type` - (Required, List: [`String`]) CDN play protocol list. Valid values: `rtmp`, `flv`, `hls`, `dash`, `hls|dash`, `customization`.
* `domain_name` - (Required, String, ForceNew) Domain name.
* `origin_address_type` - (Required, Int) Origin address type. `1`: IP. `2`: domain name.
* `origin_address` - (Required, List: [`String`]) Origin address list. Each item format: `host:port`. Port can be empty but colon is required.
* `origin_stream_play_type` - (Required, String) Origin stream play protocol. Valid values: `rtmp`, `flv`, `hls`, `dash`, `hls|dash`, `customization`.
* `origin_stream_type` - (Required, Int) Origin type. `1`: live origin. `2`: mediaPackage.
* `cache_follow_origin` - (Optional, String) Follow origin cache. Valid values: `on`, `off`. Effective only when `origin_stream_play_type` is `hls`.
* `cache_format_rule` - (Optional, Int) Cache format rule. `0`: default. `1`: live origin format. Effective only when `origin_stream_play_type` is `customization`.
* `cache_status_code` - (Optional, List: [`String`]) Status code cache list. Format: `cacheKey:interval`. Effective only when `origin_stream_play_type` is `hls`.
* `customer_name` - (Optional, String) Custom name.
* `customization_rules` - (Optional, List) Customization rules list. Effective only when `origin_stream_play_type` is `customization`.
* `follow_redirect` - (Optional, String) Follow 301/302. Valid values: `on`, `off`. Effective only when `origin_stream_play_type` is `hls`.
* `fragment_cache` - (Optional, Int) Fragment cache in ms, range: 1~60000, default: 10000. Effective only when `origin_stream_play_type` is `hls`.
* `fragment_header` - (Optional, List: [`String`]) Fragment custom headers, max 10 items. Effective only when `origin_stream_play_type` is `hls`.
* `fragment_keep_param` - (Optional, List: [`String`]) Fragment cache keep param list, max 30 items. Effective only when `origin_stream_play_type` is `hls`.
* `hls_play_fragment_count` - (Optional, Int) Fragment count, range: 1~10, default: 3.
* `hls_play_fragment_duration` - (Optional, Int) Fragment duration in ms, range: 1~10000, default: 3000.
* `indexer_cache` - (Optional, Int) Index cache in ms, range: 1~60000, default: 10000. Effective only when `origin_stream_play_type` is `hls`.
* `indexer_header` - (Optional, List: [`String`]) Index custom headers, max 10 items. Effective only when `origin_stream_play_type` is `hls`.
* `indexer_keep_param` - (Optional, List: [`String`]) Index cache keep param list, max 30 items. Effective only when `origin_stream_play_type` is `hls`.
* `media_package_channel_types` - (Optional, List: [`String`]) MediaPackage channel types. Valid values: `normal`, `ssai`, `linear_assembly`. Effective only when `origin_stream_type` is `2` and `media_package_type` is `media_package`.
* `media_package_type` - (Optional, String) MediaPackage type. Valid values: `media_package`, `media_package_pure_ad`, `media_package_mix_ad`. Effective only when `origin_stream_type` is `2`.
* `options_request` - (Optional, String) OPTIONS support. Valid values: `on`, `off`. Effective only when `origin_stream_play_type` is `hls`.
* `origin_host` - (Optional, String) Origin host. Effective only when `origin_stream_play_type` is `hls`.
* `origin_retry_times` - (Optional, Int) Retry count, range: 1~10, default: 10.
* `origin_timeout` - (Optional, Int) Timeout in ms, range: 1~60000, default: 10000.
* `pass_through_http_header` - (Optional, String) Whether to pass through HTTP headers. Valid values: `on`, `off`. Effective only when `origin_stream_play_type` is `hls`.
* `pass_through_param` - (Optional, String) Whether to pass through parameters. Valid values: `on`, `off`. Effective only when `origin_stream_play_type` is `hls`.
* `pass_through_response` - (Optional, String) Whether to pass through response. Valid values: `on`, `off`. Effective only when `origin_stream_play_type` is `hls`.
* `time_jitter` - (Optional, String) Timestamp correction. Valid values: `on`, `off`. Effective only when `origin_stream_play_type` is `rtmp` or `flv`.
* `url_replace_rules` - (Optional, List: [`String`]) URL rewrite rules. Format: `url1<|>url2`. Effective only when `origin_stream_play_type` is `hls`.
* `using_https` - (Optional, String) HTTPS back-to-origin. Valid values: `on`, `off`. Effective only when `origin_stream_play_type` is `flv` or `hls`.

The `customization_rules` object supports the following:

* `match_rule` - (Required, String) Match rule. Valid values: `.m3u8`, `.mpd`, `.ts`, `.mp4`, `.m4s`, `.m4a`, `.m4i`, `.m4v`, `.m4f`, `.aac`, `.webm`.
* `origin_address_type` - (Required, Int) Origin address type. `1`: IP. `2`: domain name.
* `origin_address` - (Required, List) Origin address list.
* `cache_status_code` - (Optional, List) Status code cache list.
* `cache` - (Optional, Int) Cache duration in s, range: 0~31536000.
* `customization_cache_follow_origin` - (Optional, Int) Custom cache follow origin. `0`: disabled. `1`: enabled.
* `http_header` - (Optional, List) Custom headers list.
* `keep_http_header` - (Optional, List) Cache HTTP header key list.
* `keep_param` - (Optional, List) Cache key list.
* `options_request` - (Optional, String) OPTIONS support. Valid values: `on`, `off`.
* `origin_host` - (Optional, String) Origin host.
* `origin_retry_times` - (Optional, Int) Retry count, range: 1~10.
* `origin_timeout` - (Optional, Int) Back-to-origin timeout in ms, range: 1~60000, default: 10000.
* `pass_through_http_header` - (Optional, String) Whether to pass through HTTP headers. Valid values: `on`, `off`.
* `pass_through_param` - (Optional, String) Whether to pass through parameters. Valid values: `on`, `off`.
* `pass_through_response` - (Optional, String) Whether to pass through response. Valid values: `on`, `off`.
* `url_replace_rules` - (Optional, List) URL rewrite rules.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `status` - Configuration status. `0`: configuring. `1`: success. `2`: closing. `3`: closed successfully.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `30m`) Used when creating the resource.
* `update` - (Defaults to `30m`) Used when updating the resource.
* `delete` - (Defaults to `30m`) Used when deleting the resource.

## Import

CSS origin stream info can be imported using the domain name, e.g.

```
terraform import tencentcloud_css_origin_stream_info.example www.demo.com
```

