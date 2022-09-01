/*
Provides a resource to create a teo zone_setting

Example Usage

```hcl
resource "tencentcloud_teo_zone_setting" "zone_setting" {
  zone_id = tencentcloud_teo_zone.zone.id

  cache {
    follow_origin {
      switch = "off"
    }

    no_cache {
      switch = "off"
    }
  }

  cache_key {
    full_url_cache = "off"
    ignore_case    = "on"

    query_string {
      action = "excludeCustom"
      switch = "on"
      value  = ["test", "apple"]
    }
  }

  cache_prefresh {
    percent = 90
    switch  = "off"
  }

  client_ip_header {
    switch = "off"
  }

  compression {
    switch = "off"
  }

  force_redirect {
    redirect_status_code = 302
    switch               = "on"
  }

  https {
    http2         = "on"
    ocsp_stapling = "off"
    tls_version   = [
      "TLSv1.2",
      "TLSv1.3",
    ]

    hsts {
      include_sub_domains = "off"
      max_age             = 0
      preload             = "off"
      switch              = "off"
    }
  }

  max_age {
    follow_origin = "off"
    max_age_time  = 600
  }

  offline_cache {
    switch = "off"
  }

  origin {
    origin_pull_protocol = "follow"
  }

  post_max_size {
    max_size = 524288000
    switch   = "on"
  }

  quic {
    switch = "on"
  }

  smart_routing {
    switch = "on"
  }

  upstream_http2 {
    switch = "off"
  }

  web_socket {
    switch  = "off"
    timeout = 30
  }
}

```
Import

teo zone_setting can be imported using the id, e.g.
```
$ terraform import tencentcloud_teo_zone_setting.zone_setting zone_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220106"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTeoZoneSetting() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTeoZoneSettingRead,
		Create: resourceTencentCloudTeoZoneSettingCreate,
		Update: resourceTencentCloudTeoZoneSettingUpdate,
		Delete: resourceTencentCloudTeoZoneSettingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Site ID.",
			},

			"cache": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "Cache expiration time configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cache": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "Cache configuration.Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Cache configuration switch.- on: Enable.- off: Disable.Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"cache_time": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Cache expiration time settings.Unit: second. The maximum value is 365 days.Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"ignore_cache_control": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies whether to enable force cache.- on: Enable.- off: Disable.Note: This field may return null, indicating that no valid value can be obtained.",
									},
								},
							},
						},
						"no_cache": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "No-cache configuration.Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Whether to cache the configuration.- on: Do not cache.- off: Cache.Note: This field may return null, indicating that no valid value can be obtained.",
									},
								},
							},
						},
						"follow_origin": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "Follows the origin server configuration.Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "According to the configuration switch of the origin site, the values are: on: open; off: off.",
									},
								},
							},
						},
					},
				},
			},

			"cache_key": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "Node cache key configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"full_url_cache": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies whether to enable full-path cache.- on: Enable full-path cache (i.e., disable Ignore Query String).- off: Disable full-path cache (i.e., enable Ignore Query String).Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"ignore_case": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies whether the cache key is case-sensitive.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"query_string": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Request parameter contained in CacheKey.Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Whether to use QueryString as part of CacheKey.- on: Enable.- off: Disable.Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"action": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "- includeCustom: Include the specified query strings.- excludeCustom: Exclude the specified query strings.Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"value": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "Array of query strings used/excluded.Note: This field may return null, indicating that no valid value can be obtained.",
									},
								},
							},
						},
					},
				},
			},

			"max_age": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "Browser cache configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_age_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the max age of the cache (in seconds). The maximum value is 365 days.Note: the value 0 means not to cache.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"follow_origin": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies whether to follow the max cache age of the origin server.- on: Enable.- off: Disable.If it&#39;s on, MaxAgeTime is ignored.Note: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},

			"offline_cache": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "Offline cache configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Whether to enable offline cache.- on: Enable.- off: Disable.Note: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},

			"quic": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "QUIC access configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Whether to enable QUIC.- on: Enable.- off: Disable.",
						},
					},
				},
			},

			"post_max_size": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "Maximum size of files transferred over POST request.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies whether to enable custom setting of the maximum file size.- on: Enable. You can set a custom max size.- off: Disable. In this case, the max size defaults to 32 MB.",
						},
						"max_size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Maximum size. Value range: 1-500 MB.Note: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},

			"compression": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "Smart compression configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Whether to enable Smart compression.- on: Enable.- off: Disable.",
						},
					},
				},
			},

			"upstream_http2": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "HTTP2 origin-pull configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Whether to enable HTTP2 origin-pull.- on: Enable.- off: Disable.",
						},
					},
				},
			},

			"force_redirect": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "Force HTTPS redirect configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Whether to enable force redirect.- on: Enable.- off: Disable.",
						},
						"redirect_status_code": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Redirection status code.- 301- 302Note: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},

			"https": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "HTTPS acceleration configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"http2": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "HTTP2 configuration switch.- on: Enable.- off: Disable.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"ocsp_stapling": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "OCSP configuration switch.- on: Enable.- off: Disable.It is disabled by default. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"tls_version": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "TLS version settings. Valid values: `TLSv1`, `TLSV1.1`, `TLSV1.2`, and `TLSv1.3`.Only consecutive versions can be enabled at the same time.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"hsts": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "HSTS Configuration.Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "- on: Enable.- off: Disable.",
									},
									"max_age": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "MaxAge value.Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"include_sub_domains": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies whether to include subdomain names. Valid values: `on` and `off`.Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"preload": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies whether to preload. Valid values: `on` and `off`.Note: This field may return null, indicating that no valid value can be obtained.",
									},
								},
							},
						},
					},
				},
			},

			"origin": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "Origin server configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"origin_pull_protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Origin-pull protocol.- http: Switch HTTPS requests to HTTP.- follow: Follow the protocol of the request.- https: Switch HTTP requests to HTTPS. This only supports port 443 on the origin server.Note: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},

			"smart_routing": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "Smart acceleration configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Whether to enable smart acceleration.- on: Enable.- off: Disable.",
						},
					},
				},
			},

			"web_socket": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "WebSocket configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Whether to enable custom WebSocket timeout setting. When it&#39;s off: it means to keep the default WebSocket connection timeout period, which is 15 seconds. To change the timeout period, please set it to on.",
						},
						"timeout": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Sets timeout period in seconds. Maximum value: 120.",
						},
					},
				},
			},

			"client_ip_header": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "Origin-pull client IP header configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies whether to enable client IP header.- on: Enable.- off: Disable.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"header_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the origin-pull client IP request header.Note: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},

			"cache_prefresh": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "Cache pre-refresh configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies whether to enable cache prefresh.- on: Enable.- off: Disable.",
						},
						"percent": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Percentage of cache time before try to prefresh. Value range: 1-99.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTeoZoneSettingCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_zone_setting.create")()
	defer inconsistentCheck(d, meta)()

	var (
		zoneId string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	d.SetId(zoneId)
	return resourceTencentCloudTeoZoneSettingUpdate(d, meta)
}

func resourceTencentCloudTeoZoneSettingRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_zone_setting.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	zoneId := d.Id()

	zoneSetting, err := service.DescribeTeoZoneSetting(ctx, zoneId)

	if err != nil {
		return err
	}

	if zoneSetting == nil {
		d.SetId("")
		return fmt.Errorf("resource `zoneSetting` %s does not exist", zoneId)
	}

	if zoneSetting.ZoneId != nil {
		_ = d.Set("zone_id", zoneSetting.ZoneId)
	}

	if zoneSetting.Cache != nil {
		cacheMap := map[string]interface{}{}
		if zoneSetting.Cache.Cache != nil {
			cacheMap := map[string]interface{}{}
			if zoneSetting.Cache.Cache.Switch != nil {
				cacheMap["switch"] = zoneSetting.Cache.Cache.Switch
			}
			if zoneSetting.Cache.Cache.CacheTime != nil {
				cacheMap["cache_time"] = zoneSetting.Cache.Cache.CacheTime
			}
			if zoneSetting.Cache.Cache.IgnoreCacheControl != nil {
				cacheMap["ignore_cache_control"] = zoneSetting.Cache.Cache.IgnoreCacheControl
			}

			cacheMap["cache"] = []interface{}{cacheMap}
		}
		if zoneSetting.Cache.NoCache != nil {
			noCacheMap := map[string]interface{}{}
			if zoneSetting.Cache.NoCache.Switch != nil {
				noCacheMap["switch"] = zoneSetting.Cache.NoCache.Switch
			}

			cacheMap["no_cache"] = []interface{}{noCacheMap}
		}
		if zoneSetting.Cache.FollowOrigin != nil {
			followOriginMap := map[string]interface{}{}
			if zoneSetting.Cache.FollowOrigin.Switch != nil {
				followOriginMap["switch"] = zoneSetting.Cache.FollowOrigin.Switch
			}

			cacheMap["follow_origin"] = []interface{}{followOriginMap}
		}

		_ = d.Set("cache", []interface{}{cacheMap})
	}

	if zoneSetting.CacheKey != nil {
		cacheKeyMap := map[string]interface{}{}
		if zoneSetting.CacheKey.FullUrlCache != nil {
			cacheKeyMap["full_url_cache"] = zoneSetting.CacheKey.FullUrlCache
		}
		if zoneSetting.CacheKey.IgnoreCase != nil {
			cacheKeyMap["ignore_case"] = zoneSetting.CacheKey.IgnoreCase
		}
		if zoneSetting.CacheKey.QueryString != nil {
			queryStringMap := map[string]interface{}{}
			if zoneSetting.CacheKey.QueryString.Switch != nil {
				queryStringMap["switch"] = zoneSetting.CacheKey.QueryString.Switch
			}
			if zoneSetting.CacheKey.QueryString.Action != nil {
				queryStringMap["action"] = zoneSetting.CacheKey.QueryString.Action
			}
			if zoneSetting.CacheKey.QueryString.Value != nil {
				queryStringMap["value"] = zoneSetting.CacheKey.QueryString.Value
			}

			cacheKeyMap["query_string"] = []interface{}{queryStringMap}
		}

		_ = d.Set("cache_key", []interface{}{cacheKeyMap})
	}

	if zoneSetting.MaxAge != nil {
		maxAgeMap := map[string]interface{}{}
		if zoneSetting.MaxAge.MaxAgeTime != nil {
			maxAgeMap["max_age_time"] = zoneSetting.MaxAge.MaxAgeTime
		}
		if zoneSetting.MaxAge.FollowOrigin != nil {
			maxAgeMap["follow_origin"] = zoneSetting.MaxAge.FollowOrigin
		}

		_ = d.Set("max_age", []interface{}{maxAgeMap})
	}

	if zoneSetting.OfflineCache != nil {
		offlineCacheMap := map[string]interface{}{}
		if zoneSetting.OfflineCache.Switch != nil {
			offlineCacheMap["switch"] = zoneSetting.OfflineCache.Switch
		}

		_ = d.Set("offline_cache", []interface{}{offlineCacheMap})
	}

	if zoneSetting.Quic != nil {
		quicMap := map[string]interface{}{}
		if zoneSetting.Quic.Switch != nil {
			quicMap["switch"] = zoneSetting.Quic.Switch
		}

		_ = d.Set("quic", []interface{}{quicMap})
	}

	if zoneSetting.PostMaxSize != nil {
		postMaxSizeMap := map[string]interface{}{}
		if zoneSetting.PostMaxSize.Switch != nil {
			postMaxSizeMap["switch"] = zoneSetting.PostMaxSize.Switch
		}
		if zoneSetting.PostMaxSize.MaxSize != nil {
			postMaxSizeMap["max_size"] = zoneSetting.PostMaxSize.MaxSize
		}

		_ = d.Set("post_max_size", []interface{}{postMaxSizeMap})
	}

	if zoneSetting.Compression != nil {
		compressionMap := map[string]interface{}{}
		if zoneSetting.Compression.Switch != nil {
			compressionMap["switch"] = zoneSetting.Compression.Switch
		}

		_ = d.Set("compression", []interface{}{compressionMap})
	}

	if zoneSetting.UpstreamHttp2 != nil {
		upstreamHttp2Map := map[string]interface{}{}
		if zoneSetting.UpstreamHttp2.Switch != nil {
			upstreamHttp2Map["switch"] = zoneSetting.UpstreamHttp2.Switch
		}

		_ = d.Set("upstream_http2", []interface{}{upstreamHttp2Map})
	}

	if zoneSetting.ForceRedirect != nil {
		forceRedirectMap := map[string]interface{}{}
		if zoneSetting.ForceRedirect.Switch != nil {
			forceRedirectMap["switch"] = zoneSetting.ForceRedirect.Switch
		}
		if zoneSetting.ForceRedirect.RedirectStatusCode != nil {
			forceRedirectMap["redirect_status_code"] = zoneSetting.ForceRedirect.RedirectStatusCode
		}

		_ = d.Set("force_redirect", []interface{}{forceRedirectMap})
	}

	if zoneSetting.Https != nil {
		httpsMap := map[string]interface{}{}
		if zoneSetting.Https.Http2 != nil {
			httpsMap["http2"] = zoneSetting.Https.Http2
		}
		if zoneSetting.Https.OcspStapling != nil {
			httpsMap["ocsp_stapling"] = zoneSetting.Https.OcspStapling
		}
		if zoneSetting.Https.TlsVersion != nil {
			httpsMap["tls_version"] = zoneSetting.Https.TlsVersion
		}
		if zoneSetting.Https.Hsts != nil {
			hstsMap := map[string]interface{}{}
			if zoneSetting.Https.Hsts.Switch != nil {
				hstsMap["switch"] = zoneSetting.Https.Hsts.Switch
			}
			if zoneSetting.Https.Hsts.MaxAge != nil {
				hstsMap["max_age"] = zoneSetting.Https.Hsts.MaxAge
			}
			if zoneSetting.Https.Hsts.IncludeSubDomains != nil {
				hstsMap["include_sub_domains"] = zoneSetting.Https.Hsts.IncludeSubDomains
			}
			if zoneSetting.Https.Hsts.Preload != nil {
				hstsMap["preload"] = zoneSetting.Https.Hsts.Preload
			}

			httpsMap["hsts"] = []interface{}{hstsMap}
		}

		_ = d.Set("https", []interface{}{httpsMap})
	}

	if zoneSetting.Origin != nil {
		originMap := map[string]interface{}{}
		if zoneSetting.Origin.OriginPullProtocol != nil {
			originMap["origin_pull_protocol"] = zoneSetting.Origin.OriginPullProtocol
		}

		_ = d.Set("origin", []interface{}{originMap})
	}

	if zoneSetting.SmartRouting != nil {
		smartRoutingMap := map[string]interface{}{}
		if zoneSetting.SmartRouting.Switch != nil {
			smartRoutingMap["switch"] = zoneSetting.SmartRouting.Switch
		}

		_ = d.Set("smart_routing", []interface{}{smartRoutingMap})
	}

	if zoneSetting.WebSocket != nil {
		webSocketMap := map[string]interface{}{}
		if zoneSetting.WebSocket.Switch != nil {
			webSocketMap["switch"] = zoneSetting.WebSocket.Switch
		}
		if zoneSetting.WebSocket.Timeout != nil {
			webSocketMap["timeout"] = zoneSetting.WebSocket.Timeout
		}

		_ = d.Set("web_socket", []interface{}{webSocketMap})
	}

	if zoneSetting.ClientIpHeader != nil {
		clientIpHeaderMap := map[string]interface{}{}
		if zoneSetting.ClientIpHeader.Switch != nil {
			clientIpHeaderMap["switch"] = zoneSetting.ClientIpHeader.Switch
		}
		if zoneSetting.ClientIpHeader.HeaderName != nil {
			clientIpHeaderMap["header_name"] = zoneSetting.ClientIpHeader.HeaderName
		}

		_ = d.Set("client_ip_header", []interface{}{clientIpHeaderMap})
	}

	if zoneSetting.CachePrefresh != nil {
		cachePrefreshMap := map[string]interface{}{}
		if zoneSetting.CachePrefresh.Switch != nil {
			cachePrefreshMap["switch"] = zoneSetting.CachePrefresh.Switch
		}
		if zoneSetting.CachePrefresh.Percent != nil {
			cachePrefreshMap["percent"] = zoneSetting.CachePrefresh.Percent
		}

		_ = d.Set("cache_prefresh", []interface{}{cachePrefreshMap})
	}

	return nil
}

func resourceTencentCloudTeoZoneSettingUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_zone_setting.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request = teo.NewModifyZoneSettingRequest()
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "cache"); ok {
		cacheConfig := teo.CacheConfig{}
		if CacheMap, ok := helper.InterfaceToMap(dMap, "cache"); ok {
			cacheConfigCache := teo.CacheConfigCache{}
			if v, ok := CacheMap["switch"]; ok {
				cacheConfigCache.Switch = helper.String(v.(string))
			}
			if v, ok := CacheMap["cache_time"]; ok {
				cacheConfigCache.CacheTime = helper.IntInt64(v.(int))
			}
			if v, ok := CacheMap["ignore_cache_control"]; ok {
				cacheConfigCache.IgnoreCacheControl = helper.String(v.(string))
			}
			cacheConfig.Cache = &cacheConfigCache
		}
		if NoCacheMap, ok := helper.InterfaceToMap(dMap, "no_cache"); ok {
			cacheConfigNoCache := teo.CacheConfigNoCache{}
			if v, ok := NoCacheMap["switch"]; ok {
				cacheConfigNoCache.Switch = helper.String(v.(string))
			}
			cacheConfig.NoCache = &cacheConfigNoCache
		}
		if FollowOriginMap, ok := helper.InterfaceToMap(dMap, "follow_origin"); ok {
			cacheConfigFollowOrigin := teo.CacheConfigFollowOrigin{}
			if v, ok := FollowOriginMap["switch"]; ok {
				cacheConfigFollowOrigin.Switch = helper.String(v.(string))
			}
			cacheConfig.FollowOrigin = &cacheConfigFollowOrigin
		}

		request.Cache = &cacheConfig
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "cache_key"); ok {
		cacheKey := teo.CacheKey{}
		if v, ok := dMap["full_url_cache"]; ok {
			cacheKey.FullUrlCache = helper.String(v.(string))
		}
		if v, ok := dMap["ignore_case"]; ok {
			cacheKey.IgnoreCase = helper.String(v.(string))
		}
		if QueryStringMap, ok := helper.InterfaceToMap(dMap, "query_string"); ok {
			queryString := teo.QueryString{}
			if v, ok := QueryStringMap["switch"]; ok {
				queryString.Switch = helper.String(v.(string))
			}
			if v, ok := QueryStringMap["action"]; ok {
				queryString.Action = helper.String(v.(string))
			}
			if v, ok := QueryStringMap["value"]; ok {
				valueSet := v.(*schema.Set).List()
				for i := range valueSet {
					value := valueSet[i].(string)
					queryString.Value = append(queryString.Value, &value)
				}
			}
			cacheKey.QueryString = &queryString
		}

		request.CacheKey = &cacheKey
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "max_age"); ok {
		maxAge := teo.MaxAge{}
		if v, ok := dMap["max_age_time"]; ok {
			maxAge.MaxAgeTime = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["follow_origin"]; ok {
			maxAge.FollowOrigin = helper.String(v.(string))
		}

		request.MaxAge = &maxAge
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "offline_cache"); ok {
		offlineCache := teo.OfflineCache{}
		if v, ok := dMap["switch"]; ok {
			offlineCache.Switch = helper.String(v.(string))
		}

		request.OfflineCache = &offlineCache
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "quic"); ok {
		quic := teo.Quic{}
		if v, ok := dMap["switch"]; ok {
			quic.Switch = helper.String(v.(string))
		}

		request.Quic = &quic
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "post_max_size"); ok {
		postMaxSize := teo.PostMaxSize{}
		if v, ok := dMap["switch"]; ok {
			postMaxSize.Switch = helper.String(v.(string))
		}
		if v, ok := dMap["max_size"]; ok {
			postMaxSize.MaxSize = helper.IntInt64(v.(int))
		}

		request.PostMaxSize = &postMaxSize
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "compression"); ok {
		compression := teo.Compression{}
		if v, ok := dMap["switch"]; ok {
			compression.Switch = helper.String(v.(string))
		}

		request.Compression = &compression
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "upstream_http2"); ok {
		upstreamHttp2 := teo.UpstreamHttp2{}
		if v, ok := dMap["switch"]; ok {
			upstreamHttp2.Switch = helper.String(v.(string))
		}

		request.UpstreamHttp2 = &upstreamHttp2
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "force_redirect"); ok {
		forceRedirect := teo.ForceRedirect{}
		if v, ok := dMap["switch"]; ok {
			forceRedirect.Switch = helper.String(v.(string))
		}
		if v, ok := dMap["redirect_status_code"]; ok {
			forceRedirect.RedirectStatusCode = helper.IntInt64(v.(int))
		}

		request.ForceRedirect = &forceRedirect
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "https"); ok {
		https := teo.Https{}
		if v, ok := dMap["http2"]; ok {
			https.Http2 = helper.String(v.(string))
		}
		if v, ok := dMap["ocsp_stapling"]; ok {
			https.OcspStapling = helper.String(v.(string))
		}
		if v, ok := dMap["tls_version"]; ok {
			tlsVersionSet := v.(*schema.Set).List()
			for i := range tlsVersionSet {
				tlsVersion := tlsVersionSet[i].(string)
				https.TlsVersion = append(https.TlsVersion, &tlsVersion)
			}
		}
		if HstsMap, ok := helper.InterfaceToMap(dMap, "hsts"); ok {
			hsts := teo.Hsts{}
			if v, ok := HstsMap["switch"]; ok {
				hsts.Switch = helper.String(v.(string))
			}
			if v, ok := HstsMap["max_age"]; ok {
				hsts.MaxAge = helper.IntInt64(v.(int))
			}
			if v, ok := HstsMap["include_sub_domains"]; ok {
				hsts.IncludeSubDomains = helper.String(v.(string))
			}
			if v, ok := HstsMap["preload"]; ok {
				hsts.Preload = helper.String(v.(string))
			}
			https.Hsts = &hsts
		}

		request.Https = &https
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "origin"); ok {
		origin := teo.Origin{}
		if v, ok := dMap["origin_pull_protocol"]; ok {
			origin.OriginPullProtocol = helper.String(v.(string))
		}

		request.Origin = &origin
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "smart_routing"); ok {
		smartRouting := teo.SmartRouting{}
		if v, ok := dMap["switch"]; ok {
			smartRouting.Switch = helper.String(v.(string))
		}

		request.SmartRouting = &smartRouting
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "web_socket"); ok {
		webSocket := teo.WebSocket{}
		if v, ok := dMap["switch"]; ok {
			webSocket.Switch = helper.String(v.(string))
		}
		if v, ok := dMap["timeout"]; ok {
			webSocket.Timeout = helper.IntInt64(v.(int))
		}

		request.WebSocket = &webSocket
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "client_ip_header"); ok {
		clientIp := teo.ClientIp{}
		if v, ok := dMap["switch"]; ok {
			clientIp.Switch = helper.String(v.(string))
		}
		if v, ok := dMap["header_name"]; ok {
			clientIp.HeaderName = helper.String(v.(string))
		}

		request.ClientIpHeader = &clientIp
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "cache_prefresh"); ok {
		cachePrefresh := teo.CachePrefresh{}
		if v, ok := dMap["switch"]; ok {
			cachePrefresh.Switch = helper.String(v.(string))
		}
		if v, ok := dMap["percent"]; ok {
			cachePrefresh.Percent = helper.IntInt64(v.(int))
		}

		request.CachePrefresh = &cachePrefresh
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyZoneSetting(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create teo zoneSetting failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTeoZoneSettingRead(d, meta)
}

func resourceTencentCloudTeoZoneSettingDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_zone_setting.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
