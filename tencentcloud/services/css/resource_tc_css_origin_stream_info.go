package css

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCssOriginStreamInfo() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCssOriginStreamInfoCreate,
		Read:   resourceTencentCloudCssOriginStreamInfoRead,
		Update: resourceTencentCloudCssOriginStreamInfoUpdate,
		Delete: resourceTencentCloudCssOriginStreamInfoDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Domain name.",
			},

			"origin_stream_play_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Origin stream play protocol. Valid values: `rtmp`, `flv`, `hls`, `dash`, `hls|dash`, `customization`.",
			},

			"cdn_stream_play_type": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "CDN play protocol list. Valid values: `rtmp`, `flv`, `hls`, `dash`, `hls|dash`, `customization`.",
			},

			"origin_stream_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Origin type. `1`: live origin. `2`: mediaPackage.",
			},

			"origin_address": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Origin address list. Each item format: `host:port`. Port can be empty but colon is required.",
			},

			"origin_address_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Origin address type. `1`: IP. `2`: domain name.",
			},

			"customer_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Custom name.",
			},

			"origin_host": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Origin host. Effective only when `origin_stream_play_type` is `hls`.",
			},

			"origin_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Timeout in ms, range: 1~60000, default: 10000.",
			},

			"origin_retry_times": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Retry count, range: 1~10, default: 10.",
			},

			"pass_through_http_header": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Whether to pass through HTTP headers. Valid values: `on`, `off`. Effective only when `origin_stream_play_type` is `hls`.",
			},

			"pass_through_response": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Whether to pass through response. Valid values: `on`, `off`. Effective only when `origin_stream_play_type` is `hls`.",
			},

			"pass_through_param": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Whether to pass through parameters. Valid values: `on`, `off`. Effective only when `origin_stream_play_type` is `hls`.",
			},

			"indexer_cache": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Index cache in ms, range: 1~60000, default: 10000. Effective only when `origin_stream_play_type` is `hls`.",
			},

			"fragment_cache": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Fragment cache in ms, range: 1~60000, default: 10000. Effective only when `origin_stream_play_type` is `hls`.",
			},

			"hls_play_fragment_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Fragment count, range: 1~10, default: 3.",
			},

			"hls_play_fragment_duration": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Fragment duration in ms, range: 1~10000, default: 3000.",
			},

			"time_jitter": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Timestamp correction. Valid values: `on`, `off`. Effective only when `origin_stream_play_type` is `rtmp` or `flv`.",
			},

			"using_https": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "HTTPS back-to-origin. Valid values: `on`, `off`. Effective only when `origin_stream_play_type` is `flv` or `hls`.",
			},

			"cache_follow_origin": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Follow origin cache. Valid values: `on`, `off`. Effective only when `origin_stream_play_type` is `hls`.",
			},

			"cache_status_code": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Status code cache list. Format: `cacheKey:interval`. Effective only when `origin_stream_play_type` is `hls`.",
			},

			"url_replace_rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "URL rewrite rules. Format: `url1<|>url2`. Effective only when `origin_stream_play_type` is `hls`.",
			},

			"options_request": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "OPTIONS support. Valid values: `on`, `off`. Effective only when `origin_stream_play_type` is `hls`.",
			},

			"follow_redirect": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Follow 301/302. Valid values: `on`, `off`. Effective only when `origin_stream_play_type` is `hls`.",
			},

			"indexer_keep_param": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Index cache keep param list, max 30 items. Effective only when `origin_stream_play_type` is `hls`.",
			},

			"fragment_keep_param": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Fragment cache keep param list, max 30 items. Effective only when `origin_stream_play_type` is `hls`.",
			},

			"media_package_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "MediaPackage type. Valid values: `media_package`, `media_package_pure_ad`, `media_package_mix_ad`. Effective only when `origin_stream_type` is `2`.",
			},

			"media_package_channel_types": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "MediaPackage channel types. Valid values: `normal`, `ssai`, `linear_assembly`. Effective only when `origin_stream_type` is `2` and `media_package_type` is `media_package`.",
			},

			"indexer_header": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Index custom headers, max 10 items. Effective only when `origin_stream_play_type` is `hls`.",
			},

			"fragment_header": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Fragment custom headers, max 10 items. Effective only when `origin_stream_play_type` is `hls`.",
			},

			"customization_rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Customization rules list. Effective only when `origin_stream_play_type` is `customization`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"match_rule": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Match rule. Valid values: `.m3u8`, `.mpd`, `.ts`, `.mp4`, `.m4s`, `.m4a`, `.m4i`, `.m4v`, `.m4f`, `.aac`, `.webm`.",
						},
						"origin_address_type": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Origin address type. `1`: IP. `2`: domain name.",
						},
						"origin_address": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Origin address list.",
						},
						"origin_host": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Origin host.",
						},
						"pass_through_http_header": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Whether to pass through HTTP headers. Valid values: `on`, `off`.",
						},
						"pass_through_response": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Whether to pass through response. Valid values: `on`, `off`.",
						},
						"pass_through_param": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Whether to pass through parameters. Valid values: `on`, `off`.",
						},
						"url_replace_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "URL rewrite rules.",
						},
						"options_request": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "OPTIONS support. Valid values: `on`, `off`.",
						},
						"origin_timeout": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Back-to-origin timeout in ms, range: 1~60000, default: 10000.",
						},
						"origin_retry_times": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Retry count, range: 1~10.",
						},
						"cache_status_code": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Status code cache list.",
						},
						"cache": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Cache duration in s, range: 0~31536000.",
						},
						"keep_param": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Cache key list.",
						},
						"http_header": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Custom headers list.",
						},
						"customization_cache_follow_origin": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Custom cache follow origin. `0`: disabled. `1`: enabled.",
						},
						"keep_http_header": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Cache HTTP header key list.",
						},
					},
				},
			},

			"cache_format_rule": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Cache format rule. `0`: default. `1`: live origin format. Effective only when `origin_stream_play_type` is `customization`.",
			},

			// computed
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Configuration status. `0`: configuring. `1`: success. `2`: closing. `3`: closed successfully.",
			},
		},
	}
}

func resourceTencentCloudCssOriginStreamInfoCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_css_origin_stream_info.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	domainName := d.Get("domain_name").(string)
	d.SetId(domainName)
	return resourceTencentCloudCssOriginStreamInfoUpdate(d, meta)
}

func resourceTencentCloudCssOriginStreamInfoRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_css_origin_stream_info.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = CssService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		domainName = d.Id()
	)

	respData, err := service.DescribeCssOriginStreamInfo(ctx, domainName)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_css_origin_stream_info` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("domain_name", domainName)

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	if respData.OriginStreamPlayType != nil {
		_ = d.Set("origin_stream_play_type", respData.OriginStreamPlayType)
	}

	if respData.CdnStreamPlayType != nil {
		_ = d.Set("cdn_stream_play_type", respData.CdnStreamPlayType)
	}

	if respData.OriginStreamType != nil {
		_ = d.Set("origin_stream_type", respData.OriginStreamType)
	}

	if respData.OriginAddress != nil {
		_ = d.Set("origin_address", respData.OriginAddress)
	}

	if respData.OriginAddressType != nil {
		_ = d.Set("origin_address_type", respData.OriginAddressType)
	}

	if respData.CustomerName != nil {
		_ = d.Set("customer_name", respData.CustomerName)
	}

	if respData.OriginHost != nil {
		_ = d.Set("origin_host", respData.OriginHost)
	}

	if respData.OriginTimeout != nil {
		_ = d.Set("origin_timeout", respData.OriginTimeout)
	}

	if respData.OriginRetryTimes != nil {
		_ = d.Set("origin_retry_times", respData.OriginRetryTimes)
	}

	if respData.PassThroughHttpHeader != nil {
		_ = d.Set("pass_through_http_header", respData.PassThroughHttpHeader)
	}

	if respData.PassThroughResponse != nil {
		_ = d.Set("pass_through_response", respData.PassThroughResponse)
	}

	if respData.PassThroughParam != nil {
		_ = d.Set("pass_through_param", respData.PassThroughParam)
	}

	if respData.IndexerCache != nil {
		_ = d.Set("indexer_cache", respData.IndexerCache)
	}

	if respData.FragmentCache != nil {
		_ = d.Set("fragment_cache", respData.FragmentCache)
	}

	if respData.HlsPlayFragmentCount != nil {
		_ = d.Set("hls_play_fragment_count", respData.HlsPlayFragmentCount)
	}

	if respData.HlsPlayFragmentDuration != nil {
		_ = d.Set("hls_play_fragment_duration", respData.HlsPlayFragmentDuration)
	}

	if respData.TimeJitter != nil {
		_ = d.Set("time_jitter", respData.TimeJitter)
	}

	if respData.UsingHttps != nil {
		_ = d.Set("using_https", respData.UsingHttps)
	}

	if respData.CacheFollowOrigin != nil {
		_ = d.Set("cache_follow_origin", respData.CacheFollowOrigin)
	}

	if respData.CacheStatusCode != nil {
		_ = d.Set("cache_status_code", respData.CacheStatusCode)
	}

	if respData.UrlReplaceRules != nil {
		_ = d.Set("url_replace_rules", respData.UrlReplaceRules)
	}

	if respData.OptionsRequest != nil {
		_ = d.Set("options_request", respData.OptionsRequest)
	}

	if respData.FollowRedirect != nil {
		_ = d.Set("follow_redirect", respData.FollowRedirect)
	}

	if respData.IndexerKeepParam != nil {
		_ = d.Set("indexer_keep_param", respData.IndexerKeepParam)
	}

	if respData.FragmentKeepParam != nil {
		_ = d.Set("fragment_keep_param", respData.FragmentKeepParam)
	}

	if respData.MediaPackageType != nil {
		_ = d.Set("media_package_type", respData.MediaPackageType)
	}

	_ = d.Set("media_package_channel_types", helper.PStrings(respData.MediaPackageChannelTypes))

	if respData.IndexerHeader != nil {
		_ = d.Set("indexer_header", respData.IndexerHeader)
	}

	if respData.FragmentHeader != nil {
		_ = d.Set("fragment_header", respData.FragmentHeader)
	}

	if respData.CustomizationRules != nil {
		_ = d.Set("customization_rules", flattenOriginStreamCustomizationRules(respData.CustomizationRules))
	}

	if respData.CacheFormatRule != nil {
		_ = d.Set("cache_format_rule", respData.CacheFormatRule)
	}

	return nil
}

func resourceTencentCloudCssOriginStreamInfoUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_css_origin_stream_info.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		domainName = d.Id()
		request    = css.NewModifyOriginStreamInfoRequest()
	)

	request.DomainName = helper.String(domainName)

	if v, ok := d.GetOk("origin_stream_play_type"); ok {
		request.OriginStreamPlayType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cdn_stream_play_type"); ok {
		request.CdnStreamPlayType = helper.InterfacesStringsPoint(v.([]interface{}))
	}

	if v, ok := d.GetOkExists("origin_stream_type"); ok {
		request.OriginStreamType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("origin_address"); ok {
		request.OriginAddress = helper.InterfacesStringsPoint(v.([]interface{}))
	}

	if v, ok := d.GetOkExists("origin_address_type"); ok {
		request.OriginAddressType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("customer_name"); ok {
		request.CustomerName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("origin_host"); ok {
		request.OriginHost = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("origin_timeout"); ok {
		request.OriginTimeout = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("origin_retry_times"); ok {
		request.OriginRetryTimes = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("pass_through_http_header"); ok {
		request.PassThroughHttpHeader = helper.String(v.(string))
	}

	if v, ok := d.GetOk("pass_through_response"); ok {
		request.PassThroughResponse = helper.String(v.(string))
	}

	if v, ok := d.GetOk("pass_through_param"); ok {
		request.PassThroughParam = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("indexer_cache"); ok {
		request.IndexerCache = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("fragment_cache"); ok {
		request.FragmentCache = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("hls_play_fragment_count"); ok {
		request.HlsPlayFragmentCount = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("hls_play_fragment_duration"); ok {
		request.HlsPlayFragmentDuration = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("time_jitter"); ok {
		request.TimeJitter = helper.String(v.(string))
	}

	if v, ok := d.GetOk("using_https"); ok {
		request.UsingHttps = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cache_follow_origin"); ok {
		request.CacheFollowOrigin = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cache_status_code"); ok {
		request.CacheStatusCode = helper.InterfacesStringsPoint(v.([]interface{}))
	}

	if v, ok := d.GetOk("url_replace_rules"); ok {
		request.UrlReplaceRules = helper.InterfacesStringsPoint(v.([]interface{}))
	}

	if v, ok := d.GetOk("options_request"); ok {
		request.OptionsRequest = helper.String(v.(string))
	}

	if v, ok := d.GetOk("follow_redirect"); ok {
		request.FollowRedirect = helper.String(v.(string))
	}

	if v, ok := d.GetOk("indexer_keep_param"); ok {
		request.IndexerKeepParam = helper.InterfacesStringsPoint(v.([]interface{}))
	}

	if v, ok := d.GetOk("fragment_keep_param"); ok {
		request.FragmentKeepParam = helper.InterfacesStringsPoint(v.([]interface{}))
	}

	if v, ok := d.GetOk("media_package_type"); ok {
		request.MediaPackageType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("media_package_channel_types"); ok {
		request.MediaPackageChannelTypes = helper.InterfacesStringsPoint(v.([]interface{}))
	}

	if v, ok := d.GetOk("indexer_header"); ok {
		request.IndexerHeader = helper.InterfacesStringsPoint(v.([]interface{}))
	}

	if v, ok := d.GetOk("fragment_header"); ok {
		request.FragmentHeader = helper.InterfacesStringsPoint(v.([]interface{}))
	}

	if v, ok := d.GetOk("customization_rules"); ok {
		request.CustomizationRules = buildOriginStreamCustomizationRules(v.([]interface{}))
	}

	if v, ok := d.GetOkExists("cache_format_rule"); ok {
		request.CacheFormatRule = helper.IntInt64(v.(int))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCssClient().ModifyOriginStreamInfoWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update css origin stream info failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// Use create timeout for new resources, update timeout for existing ones.
	pollTimeout := d.Timeout(schema.TimeoutUpdate)
	if d.IsNewResource() {
		pollTimeout = d.Timeout(schema.TimeoutCreate)
	}

	// Poll until Status == 1 (success)
	service := CssService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	pollErr := resource.RetryContext(ctx, pollTimeout, func() *resource.RetryError {
		respData, e := service.DescribeCssOriginStreamInfo(ctx, domainName)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if respData == nil || respData.Status == nil {
			return resource.RetryableError(fmt.Errorf("waiting for css origin stream info [%s] to be ready", domainName))
		}

		status := *respData.Status
		if status == int64(1) {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("css origin stream info [%s] is not ready, current status: %d", domainName, status))
	})

	if pollErr != nil {
		log.Printf("[CRITAL]%s polling css origin stream info status failed, reason:%+v", logId, pollErr)
		return pollErr
	}

	return resourceTencentCloudCssOriginStreamInfoRead(d, meta)
}

func resourceTencentCloudCssOriginStreamInfoDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_css_origin_stream_info.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		domainName = d.Id()
		request    = css.NewCloseSourceStreamRequest()
	)

	request.DomainName = helper.String(domainName)

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCssClient().CloseSourceStreamWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete css origin stream info failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// Poll until Status == 3 (closed successfully)
	service := CssService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	pollErr := resource.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		respData, e := service.DescribeCssOriginStreamInfo(ctx, domainName)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if respData == nil || respData.Status == nil {
			return nil
		}

		status := *respData.Status
		if status == int64(3) {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("css origin stream info [%s] is not closed yet, current status: %d", domainName, status))
	})

	if pollErr != nil {
		log.Printf("[CRITAL]%s polling css origin stream info close status failed, reason:%+v", logId, pollErr)
		return pollErr
	}

	return nil
}

// buildOriginStreamCustomizationRules converts terraform list to SDK struct slice.
func buildOriginStreamCustomizationRules(rawList []interface{}) []*css.OriginStreamCustomizationRule {
	rules := make([]*css.OriginStreamCustomizationRule, 0, len(rawList))
	for _, item := range rawList {
		m := item.(map[string]interface{})
		rule := &css.OriginStreamCustomizationRule{}

		if v, ok := m["match_rule"].(string); ok && v != "" {
			rule.MatchRule = helper.String(v)
		}

		if v, ok := m["origin_address_type"].(int); ok {
			rule.OriginAddressType = helper.IntInt64(v)
		}

		if v, ok := m["origin_address"]; ok {
			rule.OriginAddress = helper.InterfacesStringsPoint(v.([]interface{}))
		}

		if v, ok := m["origin_host"].(string); ok && v != "" {
			rule.OriginHost = helper.String(v)
		}

		if v, ok := m["pass_through_http_header"].(string); ok && v != "" {
			rule.PassThroughHttpHeader = helper.String(v)
		}

		if v, ok := m["pass_through_response"].(string); ok && v != "" {
			rule.PassThroughResponse = helper.String(v)
		}

		if v, ok := m["pass_through_param"].(string); ok && v != "" {
			rule.PassThroughParam = helper.String(v)
		}

		if v, ok := m["url_replace_rules"]; ok {
			rule.UrlReplaceRules = helper.InterfacesStringsPoint(v.([]interface{}))
		}

		if v, ok := m["options_request"].(string); ok && v != "" {
			rule.OptionsRequest = helper.String(v)
		}

		if v, ok := m["origin_timeout"].(int); ok && v != 0 {
			rule.OriginTimeout = helper.IntInt64(v)
		}

		if v, ok := m["origin_retry_times"].(int); ok && v != 0 {
			rule.OriginRetryTimes = helper.IntInt64(v)
		}

		if v, ok := m["cache_status_code"]; ok {
			rule.CacheStatusCode = helper.InterfacesStringsPoint(v.([]interface{}))
		}

		if v, ok := m["cache"].(int); ok && v != 0 {
			rule.Cache = helper.IntInt64(v)
		}

		if v, ok := m["keep_param"]; ok {
			rule.KeepParam = helper.InterfacesStringsPoint(v.([]interface{}))
		}

		if v, ok := m["http_header"]; ok {
			rule.HttpHeader = helper.InterfacesStringsPoint(v.([]interface{}))
		}

		if v, ok := m["customization_cache_follow_origin"].(int); ok {
			rule.CustomizationCacheFollowOrigin = helper.IntInt64(v)
		}

		if v, ok := m["keep_http_header"]; ok {
			rule.KeepHttpHeader = helper.InterfacesStringsPoint(v.([]interface{}))
		}

		rules = append(rules, rule)
	}

	return rules
}

// flattenOriginStreamCustomizationRules converts SDK struct slice to terraform list.
func flattenOriginStreamCustomizationRules(rules []*css.OriginStreamCustomizationRule) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(rules))
	for _, rule := range rules {
		m := map[string]interface{}{}

		if rule.MatchRule != nil {
			m["match_rule"] = *rule.MatchRule
		}

		if rule.OriginAddressType != nil {
			m["origin_address_type"] = int(*rule.OriginAddressType)
		}

		if rule.OriginAddress != nil {
			m["origin_address"] = rule.OriginAddress
		}

		if rule.OriginHost != nil {
			m["origin_host"] = *rule.OriginHost
		}

		if rule.PassThroughHttpHeader != nil {
			m["pass_through_http_header"] = *rule.PassThroughHttpHeader
		}

		if rule.PassThroughResponse != nil {
			m["pass_through_response"] = *rule.PassThroughResponse
		}

		if rule.PassThroughParam != nil {
			m["pass_through_param"] = *rule.PassThroughParam
		}

		if rule.UrlReplaceRules != nil {
			m["url_replace_rules"] = rule.UrlReplaceRules
		}

		if rule.OptionsRequest != nil {
			m["options_request"] = *rule.OptionsRequest
		}

		if rule.OriginTimeout != nil {
			m["origin_timeout"] = int(*rule.OriginTimeout)
		}

		if rule.OriginRetryTimes != nil {
			m["origin_retry_times"] = int(*rule.OriginRetryTimes)
		}

		if rule.CacheStatusCode != nil {
			m["cache_status_code"] = rule.CacheStatusCode
		}

		if rule.Cache != nil {
			m["cache"] = int(*rule.Cache)
		}

		if rule.KeepParam != nil {
			m["keep_param"] = rule.KeepParam
		}

		if rule.HttpHeader != nil {
			m["http_header"] = rule.HttpHeader
		}

		if rule.CustomizationCacheFollowOrigin != nil {
			m["customization_cache_follow_origin"] = int(*rule.CustomizationCacheFollowOrigin)
		}

		if rule.KeepHttpHeader != nil {
			m["keep_http_header"] = rule.KeepHttpHeader
		}

		result = append(result, m)
	}

	return result
}
