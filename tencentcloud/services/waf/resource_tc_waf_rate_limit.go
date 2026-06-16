package waf

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wafv20180125 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWafRateLimit() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafRateLimitCreate,
		Read:   resourceTencentCloudWafRateLimitRead,
		Update: resourceTencentCloudWafRateLimitUpdate,
		Delete: resourceTencentCloudWafRateLimitDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Domain name.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule name.",
			},
			"priority": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Rule priority.",
			},
			"status": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Rule switch, 0: off, 1: on.",
			},
			"limit_window": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Rate limit window configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"second": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Maximum requests allowed per second.",
						},
						"minute": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Maximum requests allowed per minute.",
						},
						"hour": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Maximum requests allowed per hour.",
						},
						"quota_share": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to share quota. Only valid when object is URL. false: URL exclusive quota, true: all URLs share quota.",
						},
					},
				},
			},
			"limit_object": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rate limit object, supports API or Domain. If based on API, LimitPaths cannot be empty.",
			},
			"limit_strategy": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Rate limit strategy, 0: observe, 1: block, 2: CAPTCHA.",
			},
			"limit_method": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Rate limit method configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"method": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Request method to rate limit.",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Match type, supports EXACT, REGEX, IN, NOT_IN, CONTAINS, NOT_CONTAINS.",
						},
					},
				},
			},
			"limit_paths": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Rate limit path configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Rate limit path.",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Match type.",
						},
					},
				},
			},
			"limit_headers": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Rate limit headers configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Header key.",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Header value.",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Match type, supports EXACT, REGEX, IN, NOT_IN, CONTAINS, NOT_CONTAINS.",
						},
					},
				},
			},
			"limit_header_name": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Rate limit based on header parameter name.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"params_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Parameter name.",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Operator, supports REGEX, IN, NOT_IN, EACH.",
						},
					},
				},
			},
			"get_params_name": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Rate limit based on GET parameter name.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"params": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Match parameter.",
						},
						"func": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Logic operator.",
						},
						"content": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Match content.",
						},
					},
				},
			},
			"get_params_value": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Rate limit based on GET parameter value.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"params": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Match parameter.",
						},
						"func": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Logic operator.",
						},
						"content": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Match content.",
						},
					},
				},
			},
			"post_params_name": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Rate limit based on POST parameter name.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"params": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Match parameter.",
						},
						"func": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Logic operator.",
						},
						"content": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Match content.",
						},
					},
				},
			},
			"post_params_value": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Rate limit based on POST parameter value.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"params": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Match parameter.",
						},
						"func": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Logic operator.",
						},
						"content": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Match content.",
						},
					},
				},
			},
			"ip_location": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Rate limit based on IP location.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"params": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Match parameter.",
						},
						"func": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Logic operator.",
						},
						"content": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Match content.",
						},
					},
				},
			},
			"redirect_info": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Redirect information. Required when LimitStrategy is redirect.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Protocol.",
						},
						"domain": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Domain.",
						},
						"url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "URL path.",
						},
					},
				},
			},
			"block_page": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Block page, 0 means 429, otherwise fill in blockPageID.",
			},
			"object_src": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Rate limit object source, 0: manual input, 1: API asset.",
			},
			"quota_share": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to share quota. Only valid when object is URL. false: URL exclusive quota, true: all URLs share quota.",
			},
			"paths_option": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Path options, can configure request method for each path.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Request path.",
						},
						"method": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Request method.",
						},
					},
				},
			},
			"order": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Rate limit execution order, 0: default, rate limit first, 1: security protection first.",
			},
			"limit_rule_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Rate limit rule ID.",
			},
		},
	}
}

func resourceTencentCloudWafRateLimitCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_rate_limit.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = wafv20180125.NewCreateRateLimitV2Request()
		response = wafv20180125.NewCreateRateLimitV2Response()
		domain   string
	)

	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
		request.Domain = helper.String(domain)
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("priority"); ok {
		request.Priority = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("status"); ok {
		request.Status = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("limit_window"); ok {
		limitWindowList := v.([]interface{})
		if len(limitWindowList) > 0 {
			limitWindowMap := limitWindowList[0].(map[string]interface{})
			limitWindow := &wafv20180125.LimitWindow{}
			if v, ok := limitWindowMap["second"].(int); ok && v != 0 {
				limitWindow.Second = helper.IntInt64(v)
			}
			if v, ok := limitWindowMap["minute"].(int); ok && v != 0 {
				limitWindow.Minute = helper.IntInt64(v)
			}
			if v, ok := limitWindowMap["hour"].(int); ok && v != 0 {
				limitWindow.Hour = helper.IntInt64(v)
			}
			if v, ok := limitWindowMap["quota_share"].(bool); ok {
				limitWindow.QuotaShare = helper.Bool(v)
			}
			request.LimitWindow = limitWindow
		}
	}

	if v, ok := d.GetOk("limit_object"); ok {
		request.LimitObject = helper.String(v.(string))
	}

	if v, ok := d.GetOk("limit_strategy"); ok {
		request.LimitStrategy = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("limit_method"); ok {
		limitMethodList := v.([]interface{})
		if len(limitMethodList) > 0 {
			limitMethodMap := limitMethodList[0].(map[string]interface{})
			limitMethod := &wafv20180125.LimitMethod{}
			if v, ok := limitMethodMap["method"].(string); ok && v != "" {
				limitMethod.Method = helper.String(v)
			}
			if v, ok := limitMethodMap["type"].(string); ok && v != "" {
				limitMethod.Type = helper.String(v)
			}
			request.LimitMethod = limitMethod
		}
	}

	if v, ok := d.GetOk("limit_paths"); ok {
		limitPathsList := v.([]interface{})
		if len(limitPathsList) > 0 {
			limitPathsMap := limitPathsList[0].(map[string]interface{})
			limitPaths := &wafv20180125.LimitPath{}
			if v, ok := limitPathsMap["path"].(string); ok && v != "" {
				limitPaths.Path = helper.String(v)
			}
			if v, ok := limitPathsMap["type"].(string); ok && v != "" {
				limitPaths.Type = helper.String(v)
			}
			request.LimitPaths = limitPaths
		}
	}

	if v, ok := d.GetOk("limit_headers"); ok {
		limitHeadersList := v.([]interface{})
		for _, item := range limitHeadersList {
			limitHeaderMap := item.(map[string]interface{})
			limitHeader := &wafv20180125.LimitHeader{}
			if v, ok := limitHeaderMap["key"].(string); ok && v != "" {
				limitHeader.Key = helper.String(v)
			}
			if v, ok := limitHeaderMap["value"].(string); ok && v != "" {
				limitHeader.Value = helper.String(v)
			}
			if v, ok := limitHeaderMap["type"].(string); ok && v != "" {
				limitHeader.Type = helper.String(v)
			}
			request.LimitHeaders = append(request.LimitHeaders, limitHeader)
		}
	}

	if v, ok := d.GetOk("limit_header_name"); ok {
		limitHeaderNameList := v.([]interface{})
		if len(limitHeaderNameList) > 0 {
			limitHeaderNameMap := limitHeaderNameList[0].(map[string]interface{})
			limitHeaderName := &wafv20180125.LimitHeaderName{}
			if v, ok := limitHeaderNameMap["params_name"].(string); ok && v != "" {
				limitHeaderName.ParamsName = helper.String(v)
			}
			if v, ok := limitHeaderNameMap["type"].(string); ok && v != "" {
				limitHeaderName.Type = helper.String(v)
			}
			request.LimitHeaderName = limitHeaderName
		}
	}

	if v, ok := d.GetOk("get_params_name"); ok {
		request.GetParamsName = buildMatchOption(v.([]interface{}))
	}

	if v, ok := d.GetOk("get_params_value"); ok {
		request.GetParamsValue = buildMatchOption(v.([]interface{}))
	}

	if v, ok := d.GetOk("post_params_name"); ok {
		request.PostParamsName = buildMatchOption(v.([]interface{}))
	}

	if v, ok := d.GetOk("post_params_value"); ok {
		request.PostParamsValue = buildMatchOption(v.([]interface{}))
	}

	if v, ok := d.GetOk("ip_location"); ok {
		request.IpLocation = buildMatchOption(v.([]interface{}))
	}

	if v, ok := d.GetOk("redirect_info"); ok {
		redirectInfoList := v.([]interface{})
		if len(redirectInfoList) > 0 {
			redirectInfoMap := redirectInfoList[0].(map[string]interface{})
			redirectInfo := &wafv20180125.RedirectInfo{}
			if v, ok := redirectInfoMap["protocol"].(string); ok && v != "" {
				redirectInfo.Protocol = helper.String(v)
			}
			if v, ok := redirectInfoMap["domain"].(string); ok && v != "" {
				redirectInfo.Domain = helper.String(v)
			}
			if v, ok := redirectInfoMap["url"].(string); ok && v != "" {
				redirectInfo.Url = helper.String(v)
			}
			request.RedirectInfo = redirectInfo
		}
	}

	if v, ok := d.GetOk("block_page"); ok {
		request.BlockPage = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("object_src"); ok {
		request.ObjectSrc = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("quota_share"); ok {
		request.QuotaShare = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("paths_option"); ok {
		pathsOptionList := v.([]interface{})
		for _, item := range pathsOptionList {
			pathItemMap := item.(map[string]interface{})
			pathItem := &wafv20180125.PathItem{}
			if v, ok := pathItemMap["path"].(string); ok && v != "" {
				pathItem.Path = helper.String(v)
			}
			if v, ok := pathItemMap["method"].(string); ok && v != "" {
				pathItem.Method = helper.String(v)
			}
			request.PathsOption = append(request.PathsOption, pathItem)
		}
	}

	if v, ok := d.GetOk("order"); ok {
		request.Order = helper.IntInt64(v.(int))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().CreateRateLimitV2WithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("create waf_rate_limit failed, response is nil"))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create waf_rate_limit failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[DEBUG]%s create waf_rate_limit, logId: %s, current d.Id(): %s", logId, logId, d.Id())

	if response.Response.LimitRuleID == nil || *response.Response.LimitRuleID == 0 {
		return fmt.Errorf("create waf_rate_limit failed, LimitRuleID is nil or zero")
	}

	limitRuleId := *response.Response.LimitRuleID
	d.SetId(strings.Join([]string{domain, strconv.FormatInt(limitRuleId, 10)}, tccommon.FILED_SP))

	return resourceTencentCloudWafRateLimitRead(d, meta)
}

func resourceTencentCloudWafRateLimitRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_rate_limit.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, id: %s", d.Id())
	}

	domain := idSplit[0]
	limitRuleIdStr := idSplit[1]
	limitRuleId, err := strconv.ParseInt(limitRuleIdStr, 10, 64)
	if err != nil {
		return fmt.Errorf("parse limit_rule_id failed, id: %s, err: %v", d.Id(), err)
	}

	var ruleData *wafv20180125.LimitRuleV2

	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		request := wafv20180125.NewDescribeRateLimitsV2Request()
		request.Domain = helper.String(domain)
		request.Id = helper.Int64(limitRuleId)

		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().DescribeRateLimitsV2WithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil || len(result.Response.RateLimits) == 0 {
			return nil
		}

		ruleData = result.Response.RateLimits[0]
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s read waf_rate_limit failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if ruleData == nil {
		log.Printf("[WARN]%s resource `waf_rate_limit` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("domain", domain)

	if ruleData.Name != nil {
		_ = d.Set("name", ruleData.Name)
	}

	if ruleData.Priority != nil {
		_ = d.Set("priority", ruleData.Priority)
	}

	if ruleData.Status != nil {
		_ = d.Set("status", ruleData.Status)
	}

	if ruleData.LimitWindow != nil {
		limitWindowMap := map[string]interface{}{}
		if ruleData.LimitWindow.Second != nil {
			limitWindowMap["second"] = *ruleData.LimitWindow.Second
		}
		if ruleData.LimitWindow.Minute != nil {
			limitWindowMap["minute"] = *ruleData.LimitWindow.Minute
		}
		if ruleData.LimitWindow.Hour != nil {
			limitWindowMap["hour"] = *ruleData.LimitWindow.Hour
		}
		if ruleData.LimitWindow.QuotaShare != nil {
			limitWindowMap["quota_share"] = *ruleData.LimitWindow.QuotaShare
		}
		_ = d.Set("limit_window", []interface{}{limitWindowMap})
	}

	if ruleData.LimitObject != nil {
		_ = d.Set("limit_object", ruleData.LimitObject)
	}

	if ruleData.LimitStrategy != nil {
		_ = d.Set("limit_strategy", ruleData.LimitStrategy)
	}

	if ruleData.LimitMethod != nil {
		limitMethodMap := map[string]interface{}{}
		if ruleData.LimitMethod.Method != nil {
			limitMethodMap["method"] = *ruleData.LimitMethod.Method
		}
		if ruleData.LimitMethod.Type != nil {
			limitMethodMap["type"] = *ruleData.LimitMethod.Type
		}
		_ = d.Set("limit_method", []interface{}{limitMethodMap})
	}

	if ruleData.LimitPaths != nil {
		limitPathsMap := map[string]interface{}{}
		if ruleData.LimitPaths.Path != nil {
			limitPathsMap["path"] = *ruleData.LimitPaths.Path
		}
		if ruleData.LimitPaths.Type != nil {
			limitPathsMap["type"] = *ruleData.LimitPaths.Type
		}
		_ = d.Set("limit_paths", []interface{}{limitPathsMap})
	}

	if ruleData.LimitHeaders != nil && len(ruleData.LimitHeaders) > 0 {
		limitHeadersList := make([]interface{}, 0, len(ruleData.LimitHeaders))
		for _, limitHeader := range ruleData.LimitHeaders {
			limitHeaderMap := map[string]interface{}{}
			if limitHeader.Key != nil {
				limitHeaderMap["key"] = *limitHeader.Key
			}
			if limitHeader.Value != nil {
				limitHeaderMap["value"] = *limitHeader.Value
			}
			if limitHeader.Type != nil {
				limitHeaderMap["type"] = *limitHeader.Type
			}
			limitHeadersList = append(limitHeadersList, limitHeaderMap)
		}
		_ = d.Set("limit_headers", limitHeadersList)
	}

	if ruleData.LimitHeaderName != nil {
		limitHeaderNameMap := map[string]interface{}{}
		if ruleData.LimitHeaderName.ParamsName != nil {
			limitHeaderNameMap["params_name"] = *ruleData.LimitHeaderName.ParamsName
		}
		if ruleData.LimitHeaderName.Type != nil {
			limitHeaderNameMap["type"] = *ruleData.LimitHeaderName.Type
		}
		_ = d.Set("limit_header_name", []interface{}{limitHeaderNameMap})
	}

	if ruleData.GetParamsName != nil {
		_ = d.Set("get_params_name", flattenMatchOption(ruleData.GetParamsName))
	}

	if ruleData.GetParamsValue != nil {
		_ = d.Set("get_params_value", flattenMatchOption(ruleData.GetParamsValue))
	}

	if ruleData.PostParamsName != nil {
		_ = d.Set("post_params_name", flattenMatchOption(ruleData.PostParamsName))
	}

	if ruleData.PostParamsValue != nil {
		_ = d.Set("post_params_value", flattenMatchOption(ruleData.PostParamsValue))
	}

	if ruleData.IpLocation != nil {
		_ = d.Set("ip_location", flattenMatchOption(ruleData.IpLocation))
	}

	if ruleData.RedirectInfo != nil {
		redirectInfoMap := map[string]interface{}{}
		if ruleData.RedirectInfo.Protocol != nil {
			redirectInfoMap["protocol"] = *ruleData.RedirectInfo.Protocol
		}
		if ruleData.RedirectInfo.Domain != nil {
			redirectInfoMap["domain"] = *ruleData.RedirectInfo.Domain
		}
		if ruleData.RedirectInfo.Url != nil {
			redirectInfoMap["url"] = *ruleData.RedirectInfo.Url
		}
		_ = d.Set("redirect_info", []interface{}{redirectInfoMap})
	}

	if ruleData.BlockPage != nil {
		_ = d.Set("block_page", ruleData.BlockPage)
	}

	if ruleData.ObjectSrc != nil {
		_ = d.Set("object_src", ruleData.ObjectSrc)
	}

	if ruleData.QuotaShare != nil {
		_ = d.Set("quota_share", ruleData.QuotaShare)
	}

	if ruleData.PathsOption != nil && len(ruleData.PathsOption) > 0 {
		pathsOptionList := make([]interface{}, 0, len(ruleData.PathsOption))
		for _, pathItem := range ruleData.PathsOption {
			pathItemMap := map[string]interface{}{}
			if pathItem.Path != nil {
				pathItemMap["path"] = *pathItem.Path
			}
			if pathItem.Method != nil {
				pathItemMap["method"] = *pathItem.Method
			}
			pathsOptionList = append(pathsOptionList, pathItemMap)
		}
		_ = d.Set("paths_option", pathsOptionList)
	}

	if ruleData.Order != nil {
		_ = d.Set("order", ruleData.Order)
	}

	if ruleData.LimitRuleID != nil {
		_ = d.Set("limit_rule_id", ruleData.LimitRuleID)
	}

	return nil
}

func resourceTencentCloudWafRateLimitUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_rate_limit.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, id: %s", d.Id())
	}

	domain := idSplit[0]
	limitRuleIdStr := idSplit[1]
	limitRuleId, err := strconv.ParseInt(limitRuleIdStr, 10, 64)
	if err != nil {
		return fmt.Errorf("parse limit_rule_id failed, id: %s, err: %v", d.Id(), err)
	}

	needChange := false
	mutableArgs := []string{
		"name", "priority", "status", "limit_window", "limit_object", "limit_strategy",
		"limit_method", "limit_paths", "limit_headers", "limit_header_name",
		"get_params_name", "get_params_value", "post_params_name", "post_params_value",
		"ip_location", "redirect_info", "block_page", "object_src", "quota_share",
		"paths_option", "order",
	}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := wafv20180125.NewUpdateRateLimitV2Request()
		request.Domain = helper.String(domain)
		request.LimitRuleId = helper.Int64(limitRuleId)

		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOk("priority"); ok {
			request.Priority = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("status"); ok {
			request.Status = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("limit_window"); ok {
			limitWindowList := v.([]interface{})
			if len(limitWindowList) > 0 {
				limitWindowMap := limitWindowList[0].(map[string]interface{})
				limitWindow := &wafv20180125.LimitWindow{}
				if v, ok := limitWindowMap["second"].(int); ok && v != 0 {
					limitWindow.Second = helper.IntInt64(v)
				}
				if v, ok := limitWindowMap["minute"].(int); ok && v != 0 {
					limitWindow.Minute = helper.IntInt64(v)
				}
				if v, ok := limitWindowMap["hour"].(int); ok && v != 0 {
					limitWindow.Hour = helper.IntInt64(v)
				}
				if v, ok := limitWindowMap["quota_share"].(bool); ok {
					limitWindow.QuotaShare = helper.Bool(v)
				}
				request.LimitWindow = limitWindow
			}
		}

		if v, ok := d.GetOk("limit_object"); ok {
			request.LimitObject = helper.String(v.(string))
		}

		if v, ok := d.GetOk("limit_strategy"); ok {
			request.LimitStrategy = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("limit_method"); ok {
			limitMethodList := v.([]interface{})
			if len(limitMethodList) > 0 {
				limitMethodMap := limitMethodList[0].(map[string]interface{})
				limitMethod := &wafv20180125.LimitMethod{}
				if v, ok := limitMethodMap["method"].(string); ok && v != "" {
					limitMethod.Method = helper.String(v)
				}
				if v, ok := limitMethodMap["type"].(string); ok && v != "" {
					limitMethod.Type = helper.String(v)
				}
				request.LimitMethod = limitMethod
			}
		}

		if v, ok := d.GetOk("limit_paths"); ok {
			limitPathsList := v.([]interface{})
			if len(limitPathsList) > 0 {
				limitPathsMap := limitPathsList[0].(map[string]interface{})
				limitPaths := &wafv20180125.LimitPath{}
				if v, ok := limitPathsMap["path"].(string); ok && v != "" {
					limitPaths.Path = helper.String(v)
				}
				if v, ok := limitPathsMap["type"].(string); ok && v != "" {
					limitPaths.Type = helper.String(v)
				}
				request.LimitPaths = limitPaths
			}
		}

		if v, ok := d.GetOk("limit_headers"); ok {
			limitHeadersList := v.([]interface{})
			for _, item := range limitHeadersList {
				limitHeaderMap := item.(map[string]interface{})
				limitHeader := &wafv20180125.LimitHeader{}
				if v, ok := limitHeaderMap["key"].(string); ok && v != "" {
					limitHeader.Key = helper.String(v)
				}
				if v, ok := limitHeaderMap["value"].(string); ok && v != "" {
					limitHeader.Value = helper.String(v)
				}
				if v, ok := limitHeaderMap["type"].(string); ok && v != "" {
					limitHeader.Type = helper.String(v)
				}
				request.LimitHeaders = append(request.LimitHeaders, limitHeader)
			}
		}

		if v, ok := d.GetOk("limit_header_name"); ok {
			limitHeaderNameList := v.([]interface{})
			if len(limitHeaderNameList) > 0 {
				limitHeaderNameMap := limitHeaderNameList[0].(map[string]interface{})
				limitHeaderName := &wafv20180125.LimitHeaderName{}
				if v, ok := limitHeaderNameMap["params_name"].(string); ok && v != "" {
					limitHeaderName.ParamsName = helper.String(v)
				}
				if v, ok := limitHeaderNameMap["type"].(string); ok && v != "" {
					limitHeaderName.Type = helper.String(v)
				}
				request.LimitHeaderName = limitHeaderName
			}
		}

		if v, ok := d.GetOk("get_params_name"); ok {
			request.GetParamsName = buildMatchOption(v.([]interface{}))
		}

		if v, ok := d.GetOk("get_params_value"); ok {
			request.GetParamsValue = buildMatchOption(v.([]interface{}))
		}

		if v, ok := d.GetOk("post_params_name"); ok {
			request.PostParamsName = buildMatchOption(v.([]interface{}))
		}

		if v, ok := d.GetOk("post_params_value"); ok {
			request.PostParamsValue = buildMatchOption(v.([]interface{}))
		}

		if v, ok := d.GetOk("ip_location"); ok {
			request.IpLocation = buildMatchOption(v.([]interface{}))
		}

		if v, ok := d.GetOk("redirect_info"); ok {
			redirectInfoList := v.([]interface{})
			if len(redirectInfoList) > 0 {
				redirectInfoMap := redirectInfoList[0].(map[string]interface{})
				redirectInfo := &wafv20180125.RedirectInfo{}
				if v, ok := redirectInfoMap["protocol"].(string); ok && v != "" {
					redirectInfo.Protocol = helper.String(v)
				}
				if v, ok := redirectInfoMap["domain"].(string); ok && v != "" {
					redirectInfo.Domain = helper.String(v)
				}
				if v, ok := redirectInfoMap["url"].(string); ok && v != "" {
					redirectInfo.Url = helper.String(v)
				}
				request.RedirectInfo = redirectInfo
			}
		}

		if v, ok := d.GetOk("block_page"); ok {
			request.BlockPage = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("object_src"); ok {
			request.ObjectSrc = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOkExists("quota_share"); ok {
			request.QuotaShare = helper.Bool(v.(bool))
		}

		if v, ok := d.GetOk("paths_option"); ok {
			pathsOptionList := v.([]interface{})
			for _, item := range pathsOptionList {
				pathItemMap := item.(map[string]interface{})
				pathItem := &wafv20180125.PathItem{}
				if v, ok := pathItemMap["path"].(string); ok && v != "" {
					pathItem.Path = helper.String(v)
				}
				if v, ok := pathItemMap["method"].(string); ok && v != "" {
					pathItem.Method = helper.String(v)
				}
				request.PathsOption = append(request.PathsOption, pathItem)
			}
		}

		if v, ok := d.GetOk("order"); ok {
			request.Order = helper.IntInt64(v.(int))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().UpdateRateLimitV2WithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update waf_rate_limit failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudWafRateLimitRead(d, meta)
}

func resourceTencentCloudWafRateLimitDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_rate_limit.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, id: %s", d.Id())
	}

	domain := idSplit[0]
	limitRuleIdStr := idSplit[1]
	limitRuleId, err := strconv.ParseInt(limitRuleIdStr, 10, 64)
	if err != nil {
		return fmt.Errorf("parse limit_rule_id failed, id: %s, err: %v", d.Id(), err)
	}

	request := wafv20180125.NewDeleteRateLimitsV2Request()
	request.Domain = helper.String(domain)
	request.LimitRuleIds = []*int64{helper.Int64(limitRuleId)}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().DeleteRateLimitsV2WithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete waf_rate_limit failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}

func buildMatchOption(list []interface{}) *wafv20180125.MatchOption {
	if len(list) == 0 {
		return nil
	}
	matchOptionMap := list[0].(map[string]interface{})
	matchOption := &wafv20180125.MatchOption{}
	if v, ok := matchOptionMap["params"].(string); ok && v != "" {
		matchOption.Params = helper.String(v)
	}
	if v, ok := matchOptionMap["func"].(string); ok && v != "" {
		matchOption.Func = helper.String(v)
	}
	if v, ok := matchOptionMap["content"].(string); ok && v != "" {
		matchOption.Content = helper.String(v)
	}
	return matchOption
}

func flattenMatchOption(matchOption *wafv20180125.MatchOption) []interface{} {
	if matchOption == nil {
		return nil
	}
	matchOptionMap := map[string]interface{}{}
	if matchOption.Params != nil {
		matchOptionMap["params"] = *matchOption.Params
	}
	if matchOption.Func != nil {
		matchOptionMap["func"] = *matchOption.Func
	}
	if matchOption.Content != nil {
		matchOptionMap["content"] = *matchOption.Content
	}
	return []interface{}{matchOptionMap}
}
