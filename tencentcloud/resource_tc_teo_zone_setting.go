/*
Provides a resource to create a teo zone_setting

Example Usage

```hcl
resource "tencentcloud_teo_zone_setting" "zone_setting" {
  zone_id = &lt;nil&gt;
    cache {
		cache {
			switch = &lt;nil&gt;
			cache_time = &lt;nil&gt;
			ignore_cache_control = &lt;nil&gt;
		}
		no_cache {
			switch = &lt;nil&gt;
		}
		follow_origin {
			switch = &lt;nil&gt;
		}

  }
  cache_key {
		full_url_cache = &lt;nil&gt;
		ignore_case = &lt;nil&gt;
		query_string {
			switch = &lt;nil&gt;
			action = &lt;nil&gt;
			value = &lt;nil&gt;
		}

  }
  max_age {
		max_age_time = &lt;nil&gt;
		follow_origin = &lt;nil&gt;

  }
  offline_cache {
		switch = &lt;nil&gt;

  }
  quic {
		switch = &lt;nil&gt;

  }
  post_max_size {
		switch = &lt;nil&gt;
		max_size = &lt;nil&gt;

  }
  compression {
		switch = &lt;nil&gt;
		algorithms = &lt;nil&gt;

  }
  upstream_http2 {
		switch = &lt;nil&gt;

  }
  force_redirect {
		switch = &lt;nil&gt;
		redirect_status_code = &lt;nil&gt;

  }
  https {
		http2 = &lt;nil&gt;
		ocsp_stapling = &lt;nil&gt;
		tls_version = &lt;nil&gt;
		hsts {
			switch = &lt;nil&gt;
			max_age = &lt;nil&gt;
			include_sub_domains = &lt;nil&gt;
			preload = &lt;nil&gt;
		}
		cert_info {
			cert_id = &lt;nil&gt;
			status = &lt;nil&gt;
		}

  }
  origin {
		origins = &lt;nil&gt;
		backup_origins = &lt;nil&gt;
		origin_pull_protocol = &lt;nil&gt;
		cos_private_access = &lt;nil&gt;

  }
  smart_routing {
		switch = &lt;nil&gt;

  }
  web_socket {
		switch = &lt;nil&gt;
		timeout = &lt;nil&gt;

  }
  client_ip_header {
		switch = &lt;nil&gt;
		header_name = &lt;nil&gt;

  }
  cache_prefresh {
		switch = &lt;nil&gt;
		percent = &lt;nil&gt;

  }
  ipv6 {
		switch = &lt;nil&gt;

  }
}
```

Import

teo zone_setting can be imported using the id, e.g.

```
terraform import tencentcloud_teo_zone_setting.zone_setting zone_setting_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudTeoZoneSetting() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoZoneSettingCreate,
		Read:   resourceTencentCloudTeoZoneSettingRead,
		Update: resourceTencentCloudTeoZoneSettingUpdate,
		Delete: resourceTencentCloudTeoZoneSettingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Site ID.",
			},

			"area": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Acceleration area of the zone. Valid values: `mainland`, `overseas`.",
			},

			"cache": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Cache expiration time configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cache": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "Cache configuration. Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Cache configuration switch.- `on`: Enable.- `off`: Disable. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"cache_time": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Cache expiration time settings.Unit: second. The maximum value is 365 days. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"ignore_cache_control": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies whether to enable force cache.- `on`: Enable.- `off`: Disable. Note: This field may return null, indicating that no valid value can be obtained.",
									},
								},
							},
						},
						"no_cache": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "No-cache configuration. Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Whether to cache the configuration.- `on`: Do not cache.- `off`: Cache. Note: This field may return null, indicating that no valid value can be obtained.",
									},
								},
							},
						},
						"follow_origin": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "Follows the origin server configuration. Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies whether to follow the origin server configuration.- `on`: Enable.- `off`: Disable. Note: This field may return null, indicating that no valid value can be obtained.",
									},
								},
							},
						},
					},
				},
			},

			"cache_key": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Node cache key configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"full_url_cache": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies whether to enable full-path cache.- `on`: Enable full-path cache (i.e., disable Ignore Query String).- `off`: Disable full-path cache (i.e., enable Ignore Query String). Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"ignore_case": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies whether the cache key is case-sensitive. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"query_string": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Request parameter contained in CacheKey. Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Whether to use QueryString as part of CacheKey.- `on`: Enable.- `off`: Disable. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"action": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "- `includeCustom`: Include the specified query strings.- `excludeCustom`: Exclude the specified query strings. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"value": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "Array of query strings used/excluded. Note: This field may return null, indicating that no valid value can be obtained.",
									},
								},
							},
						},
					},
				},
			},

			"max_age": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Browser cache configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_age_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the max age of the cache (in seconds). The maximum value is 365 days. Note: the value 0 means not to cache. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"follow_origin": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies whether to follow the max cache age of the origin server.- `on`: Enable.- `off`: Disable.If it&amp;#39;s on, MaxAgeTime is ignored. Note: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},

			"offline_cache": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Offline cache configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Whether to enable offline cache.- `on`: Enable.- `off`: Disable. Note: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},

			"quic": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "QUIC access configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Whether to enable QUIC.- `on`: Enable.- `off`: Disable.",
						},
					},
				},
			},

			"post_max_size": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Maximum size of files transferred over POST request.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies whether to enable custom setting of the maximum file size.- `on`: Enable. You can set a custom max size.- `off`: Disable. In this case, the max size defaults to 32 MB.",
						},
						"max_size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Maximum size. Value range: 1-500 MB. Note: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},

			"compression": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Smart compression configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Whether to enable Smart compression.- `on`: Enable.- `off`: Disable.",
						},
						"algorithms": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Computed:    true,
							Description: "Compression algorithms to select. Valid values: `brotli`, `gzip`.",
						},
					},
				},
			},

			"upstream_http2": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "HTTP2 origin-pull configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Whether to enable HTTP2 origin-pull.- `on`: Enable.- `off`: Disable.",
						},
					},
				},
			},

			"force_redirect": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Force HTTPS redirect configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Whether to enable force redirect.- `on`: Enable.- `off`: Disable.",
						},
						"redirect_status_code": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Redirection status code.- 301- 302 Note: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},

			"https": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "HTTPS acceleration configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"http2": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "HTTP2 configuration switch.- `on`: Enable.- `off`: Disable. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"ocsp_stapling": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "OCSP configuration switch.- `on`: Enable.- `off`: Disable.It is disabled by default. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"tls_version": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "TLS version settings. Valid values: `TLSv1`, `TLSV1.1`, `TLSV1.2`, and `TLSv1.3`.Only consecutive versions can be enabled at the same time. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"hsts": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "HSTS Configuration. Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "- `on`: Enable.- `off`: Disable.",
									},
									"max_age": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "MaxAge value in seconds, should be no more than 1 day. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"include_sub_domains": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies whether to include subdomain names. Valid values: `on` and `off`. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"preload": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies whether to preload. Valid values: `on` and `off`. Note: This field may return null, indicating that no valid value can be obtained.",
									},
								},
							},
						},
						"cert_info": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Server certificate configuration. Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cert_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Server certificate ID, which is the ID of the default certificate. If you choose to upload an external certificate for SSL certificate management, a certificate ID will be generated.",
									},
									"alias": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Alias of the certificate. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Certificate type.- `default`: Default certificate.- `upload`: External certificate.- `managed`: Tencent Cloud managed certificate. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"expire_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Time when the certificate expires. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"effective_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Time when the certificate takes effect. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"status": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Certificate deployment status.- `processing`: Deploying- `deployed`: Deployed Note: This field may return null, indicating that no valid value can be obtained.",
									},
								},
							},
						},
					},
				},
			},

			"origin": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Origin server configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"origins": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Computed:    true,
							Description: "Origin sites list. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"backup_origins": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Computed:    true,
							Description: "Backup origin sites list. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"origin_pull_protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Origin-pull protocol.- `http`: Switch HTTPS requests to HTTP.- `follow`: Follow the protocol of the request.- `https`: Switch HTTP requests to HTTPS. This only supports port 443 on the origin server. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"cos_private_access": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Whether access private cos bucket is allowed when `OriginType` is cos. Note: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},

			"smart_routing": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Smart acceleration configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Whether to enable smart acceleration.- `on`: Enable.- `off`: Disable.",
						},
					},
				},
			},

			"web_socket": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "WebSocket configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Whether to enable custom WebSocket timeout setting. When it&amp;#39;s off: it means to keep the default WebSocket connection timeout period, which is 15 seconds. To change the timeout period, please set it to on.",
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
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Origin-pull client IP header configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies whether to enable client IP header.- `on`: Enable.- `off`: Disable. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"header_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the origin-pull client IP request header. Note: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},

			"cache_prefresh": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Cache pre-refresh configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies whether to enable cache prefresh.- `on`: Enable.- `off`: Disable.",
						},
						"percent": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Percentage of cache time before try to prefresh. Valid value range: 1-99.",
						},
					},
				},
			},

			"ipv6": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "IPv6 access configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "- `on`: Enable.- `off`: Disable.",
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

	var zoneId string
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

	zoneSettingId := d.Id()

	zoneSetting, err := service.DescribeTeoZoneSettingById(ctx, zoneId)
	if err != nil {
		return err
	}

	if zoneSetting == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TeoZoneSetting` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if zoneSetting.ZoneId != nil {
		_ = d.Set("zone_id", zoneSetting.ZoneId)
	}

	if zoneSetting.Area != nil {
		_ = d.Set("area", zoneSetting.Area)
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

		if zoneSetting.Compression.Algorithms != nil {
			compressionMap["algorithms"] = zoneSetting.Compression.Algorithms
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

		if zoneSetting.Https.CertInfo != nil {
			certInfoList := []interface{}{}
			for _, certInfo := range zoneSetting.Https.CertInfo {
				certInfoMap := map[string]interface{}{}

				if certInfo.CertId != nil {
					certInfoMap["cert_id"] = certInfo.CertId
				}

				if certInfo.Alias != nil {
					certInfoMap["alias"] = certInfo.Alias
				}

				if certInfo.Type != nil {
					certInfoMap["type"] = certInfo.Type
				}

				if certInfo.ExpireTime != nil {
					certInfoMap["expire_time"] = certInfo.ExpireTime
				}

				if certInfo.EffectiveTime != nil {
					certInfoMap["effective_time"] = certInfo.EffectiveTime
				}

				if certInfo.Status != nil {
					certInfoMap["status"] = certInfo.Status
				}

				certInfoList = append(certInfoList, certInfoMap)
			}

			httpsMap["cert_info"] = []interface{}{certInfoList}
		}

		_ = d.Set("https", []interface{}{httpsMap})
	}

	if zoneSetting.Origin != nil {
		originMap := map[string]interface{}{}

		if zoneSetting.Origin.Origins != nil {
			originMap["origins"] = zoneSetting.Origin.Origins
		}

		if zoneSetting.Origin.BackupOrigins != nil {
			originMap["backup_origins"] = zoneSetting.Origin.BackupOrigins
		}

		if zoneSetting.Origin.OriginPullProtocol != nil {
			originMap["origin_pull_protocol"] = zoneSetting.Origin.OriginPullProtocol
		}

		if zoneSetting.Origin.CosPrivateAccess != nil {
			originMap["cos_private_access"] = zoneSetting.Origin.CosPrivateAccess
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

	if zoneSetting.Ipv6 != nil {
		ipv6Map := map[string]interface{}{}

		if zoneSetting.Ipv6.Switch != nil {
			ipv6Map["switch"] = zoneSetting.Ipv6.Switch
		}

		_ = d.Set("ipv6", []interface{}{ipv6Map})
	}

	return nil
}

func resourceTencentCloudTeoZoneSettingUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_zone_setting.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := teo.NewModifyZoneSettingRequest()

	zoneSettingId := d.Id()

	request.ZoneId = &zoneId

	immutableArgs := []string{"zone_id", "area", "cache", "cache_key", "max_age", "offline_cache", "quic", "post_max_size", "compression", "upstream_http2", "force_redirect", "https", "origin", "smart_routing", "web_socket", "client_ip_header", "cache_prefresh", "ipv6"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("cache") {
		if dMap, ok := helper.InterfacesHeadMap(d, "cache"); ok {
			cacheConfig := teo.CacheConfig{}
			if cacheMap, ok := helper.InterfaceToMap(dMap, "cache"); ok {
				cacheConfigCache := teo.CacheConfigCache{}
				if v, ok := cacheMap["switch"]; ok {
					cacheConfigCache.Switch = helper.String(v.(string))
				}
				if v, ok := cacheMap["cache_time"]; ok {
					cacheConfigCache.CacheTime = helper.IntInt64(v.(int))
				}
				if v, ok := cacheMap["ignore_cache_control"]; ok {
					cacheConfigCache.IgnoreCacheControl = helper.String(v.(string))
				}
				cacheConfig.Cache = &cacheConfigCache
			}
			if noCacheMap, ok := helper.InterfaceToMap(dMap, "no_cache"); ok {
				cacheConfigNoCache := teo.CacheConfigNoCache{}
				if v, ok := noCacheMap["switch"]; ok {
					cacheConfigNoCache.Switch = helper.String(v.(string))
				}
				cacheConfig.NoCache = &cacheConfigNoCache
			}
			if followOriginMap, ok := helper.InterfaceToMap(dMap, "follow_origin"); ok {
				cacheConfigFollowOrigin := teo.CacheConfigFollowOrigin{}
				if v, ok := followOriginMap["switch"]; ok {
					cacheConfigFollowOrigin.Switch = helper.String(v.(string))
				}
				cacheConfig.FollowOrigin = &cacheConfigFollowOrigin
			}
			request.Cache = &cacheConfig
		}
	}

	if d.HasChange("cache_key") {
		if dMap, ok := helper.InterfacesHeadMap(d, "cache_key"); ok {
			cacheKey := teo.CacheKey{}
			if v, ok := dMap["full_url_cache"]; ok {
				cacheKey.FullUrlCache = helper.String(v.(string))
			}
			if v, ok := dMap["ignore_case"]; ok {
				cacheKey.IgnoreCase = helper.String(v.(string))
			}
			if queryStringMap, ok := helper.InterfaceToMap(dMap, "query_string"); ok {
				queryString := teo.QueryString{}
				if v, ok := queryStringMap["switch"]; ok {
					queryString.Switch = helper.String(v.(string))
				}
				if v, ok := queryStringMap["action"]; ok {
					queryString.Action = helper.String(v.(string))
				}
				if v, ok := queryStringMap["value"]; ok {
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
	}

	if d.HasChange("max_age") {
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
	}

	if d.HasChange("offline_cache") {
		if dMap, ok := helper.InterfacesHeadMap(d, "offline_cache"); ok {
			offlineCache := teo.OfflineCache{}
			if v, ok := dMap["switch"]; ok {
				offlineCache.Switch = helper.String(v.(string))
			}
			request.OfflineCache = &offlineCache
		}
	}

	if d.HasChange("quic") {
		if dMap, ok := helper.InterfacesHeadMap(d, "quic"); ok {
			quic := teo.Quic{}
			if v, ok := dMap["switch"]; ok {
				quic.Switch = helper.String(v.(string))
			}
			request.Quic = &quic
		}
	}

	if d.HasChange("post_max_size") {
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
	}

	if d.HasChange("compression") {
		if dMap, ok := helper.InterfacesHeadMap(d, "compression"); ok {
			compression := teo.Compression{}
			if v, ok := dMap["switch"]; ok {
				compression.Switch = helper.String(v.(string))
			}
			if v, ok := dMap["algorithms"]; ok {
				algorithmsSet := v.(*schema.Set).List()
				for i := range algorithmsSet {
					algorithms := algorithmsSet[i].(string)
					compression.Algorithms = append(compression.Algorithms, &algorithms)
				}
			}
			request.Compression = &compression
		}
	}

	if d.HasChange("upstream_http2") {
		if dMap, ok := helper.InterfacesHeadMap(d, "upstream_http2"); ok {
			upstreamHttp2 := teo.UpstreamHttp2{}
			if v, ok := dMap["switch"]; ok {
				upstreamHttp2.Switch = helper.String(v.(string))
			}
			request.UpstreamHttp2 = &upstreamHttp2
		}
	}

	if d.HasChange("force_redirect") {
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
	}

	if d.HasChange("https") {
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
			if hstsMap, ok := helper.InterfaceToMap(dMap, "hsts"); ok {
				hsts := teo.Hsts{}
				if v, ok := hstsMap["switch"]; ok {
					hsts.Switch = helper.String(v.(string))
				}
				if v, ok := hstsMap["max_age"]; ok {
					hsts.MaxAge = helper.IntInt64(v.(int))
				}
				if v, ok := hstsMap["include_sub_domains"]; ok {
					hsts.IncludeSubDomains = helper.String(v.(string))
				}
				if v, ok := hstsMap["preload"]; ok {
					hsts.Preload = helper.String(v.(string))
				}
				https.Hsts = &hsts
			}
			if v, ok := dMap["cert_info"]; ok {
				for _, item := range v.([]interface{}) {
					certInfoMap := item.(map[string]interface{})
					serverCertInfo := teo.ServerCertInfo{}
					if v, ok := certInfoMap["cert_id"]; ok {
						serverCertInfo.CertId = helper.String(v.(string))
					}
					if v, ok := certInfoMap["status"]; ok {
						serverCertInfo.Status = helper.String(v.(string))
					}
					https.CertInfo = append(https.CertInfo, &serverCertInfo)
				}
			}
			request.Https = &https
		}
	}

	if d.HasChange("origin") {
		if dMap, ok := helper.InterfacesHeadMap(d, "origin"); ok {
			origin := teo.Origin{}
			if v, ok := dMap["origins"]; ok {
				originsSet := v.(*schema.Set).List()
				for i := range originsSet {
					origins := originsSet[i].(string)
					origin.Origins = append(origin.Origins, &origins)
				}
			}
			if v, ok := dMap["backup_origins"]; ok {
				backupOriginsSet := v.(*schema.Set).List()
				for i := range backupOriginsSet {
					backupOrigins := backupOriginsSet[i].(string)
					origin.BackupOrigins = append(origin.BackupOrigins, &backupOrigins)
				}
			}
			if v, ok := dMap["origin_pull_protocol"]; ok {
				origin.OriginPullProtocol = helper.String(v.(string))
			}
			if v, ok := dMap["cos_private_access"]; ok {
				origin.CosPrivateAccess = helper.String(v.(string))
			}
			request.Origin = &origin
		}
	}

	if d.HasChange("smart_routing") {
		if dMap, ok := helper.InterfacesHeadMap(d, "smart_routing"); ok {
			smartRouting := teo.SmartRouting{}
			if v, ok := dMap["switch"]; ok {
				smartRouting.Switch = helper.String(v.(string))
			}
			request.SmartRouting = &smartRouting
		}
	}

	if d.HasChange("web_socket") {
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
	}

	if d.HasChange("client_ip_header") {
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
	}

	if d.HasChange("cache_prefresh") {
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
	}

	if d.HasChange("ipv6") {
		if dMap, ok := helper.InterfacesHeadMap(d, "ipv6"); ok {
			ipv6Access := teo.Ipv6Access{}
			if v, ok := dMap["switch"]; ok {
				ipv6Access.Switch = helper.String(v.(string))
			}
			request.Ipv6 = &ipv6Access
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyZoneSetting(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update teo zoneSetting failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTeoZoneSettingRead(d, meta)
}

func resourceTencentCloudTeoZoneSettingDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_zone_setting.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
