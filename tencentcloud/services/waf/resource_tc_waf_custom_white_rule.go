package waf

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"
	wafv20180125 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWafCustomWhiteRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafCustomWhiteRuleCreate,
		Read:   resourceTencentCloudWafCustomWhiteRuleRead,
		Update: resourceTencentCloudWafCustomWhiteRuleUpdate,
		Delete: resourceTencentCloudWafCustomWhiteRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Rule Name.",
			},
			"sort_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Priority, value range 1-100, The smaller the number, the higher the execution priority of this rule.",
			},
			"expire_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Expiration time in second-level timestamp, for example, 1677254399 indicates the expiration time is 2023-02-24 23:59:59; 0 indicates it will never expire.",
			},
			"strategies": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "Strategies detail.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field": {
							Type:     schema.TypeString,
							Required: true,
							Description: `Matching field
                            Different matching fields result in different matching parameters, logical operators, and matching contents. The details are as follows:
                        	<table><thead><tr><th>Matching Field</th><th>Matching Parameter</th><th>Logical Symbol</th><th>Matching Content</th></tr></thead><tbody><tr><td>IP (source IP)</td><td>Parameters are not supported.</td><td>ipmatch (match)<br>ipnmatch (mismatch)</td><td>Multiple IP addresses are separated by commas. A maximum of 20 IP addresses are allowed.</td></tr><tr><td>IPv6 (source IPv6)</td><td>Parameters are not supported.</td><td>ipmatch (match)<br>ipnmatch (mismatch)</td><td>A single IPv6 address is supported.</td></tr><tr><td>Referer (referer)</td><td>Parameters are not supported.</td><td>empty (Content is empty.)<br>null (do not exist)<br>eq (equal to)<br>neq (not equal to)<br>contains (contain)<br>ncontains (do not contain)<br>len_eq (length equals to)<br>len_gt (length is greater than)<br>len_lt (length is less than)<br>strprefix (prefix matching)<br>strsuffix (suffix matching)<br>rematch (regular expression matching)</td><td>Enter the content, with a maximum of 512 characters.</td></tr><tr><td>URL (request path)</td><td>Parameters are not supported.</td><td>eq (equal to)<br>neq (not equal to)<br>contains (contain)<br>ncontains (do not contain)<br>len_eq (length equals to)<br>len_gt (length is greater than)<br>len_lt (length is less than)<br>strprefix (prefix matching)<br>strsuffix (suffix matching)<br>rematch (regular expression matching)</td><td>Enter the content starting with /, with a maximum of 512 characters.</td></tr><tr><td>UserAgent (UserAgent)</td><td>Parameters are not supported.</td><td>Same logical symbols as the matching field <font color="Red">Referer</font></td><td>Enter the content with a maximum of 512 characters.</td></tr><tr><td>HTTP_METHOD (HTTP request method)</td><td>Parameters are not supported.</td><td>eq (equal to)<br>neq (not equal to)</td><td>Enter the method name. The uppercase is recommended.</td></tr><tr><td>QUERY_STRING (request string)</td><td>Parameters are not supported.</td><td>Same logical symbol as the matching field <font color="Red">Request Path</font></td><td>Enter the content with a maximum of 512 characters.</td></tr><tr><td>GET (GET parameter value)</td><td>Parameter entry is supported.</td><td>contains (contain)<br>ncontains (do not contain)<br>len_eq (length equals to)<br>len_gt (length is greater than)<br>len_lt (length is less than)<br>strprefix (prefix matching)<br>strsuffix (suffix matching)</td><td>Enter the content with a maximum of 512 characters.</td></tr><tr><td>GET_PARAMS_NAMES (GET parameter name)</td><td>Parameters are not supported.</td><td>exist (Parameter exists.)<br>nexist (Parameter does not exist.)<br>len_eq (length equals to)<br>len_gt (length is greater than)<br>len_lt (length is less than)<br>strprefix (prefix matching)<br>strsuffix (suffix matching)</td><td>Enter the content with a maximum of 512 characters.</td></tr><tr><td>POST (POST parameter value)</td><td>Parameter entry is supported.</td><td>Same logical symbol as the matching field <font color="Red">GET Parameter Value</font></td><td>Enter the content with a maximum of 512 characters.</td></tr><tr><td>GET_POST_NAMES (POST parameter name)</td><td>Parameters are not supported.</td><td>Same logical symbol as the matching field <font color="Red">GET Parameter Name</font></td><td>Enter the content with a maximum of 512 characters.</td></tr><tr><td>POST_BODY (complete body)</td><td>Parameters are not supported.</td><td>Same logical symbol as the matching field <font color="Red">Request Path</font></td><td>Enter the body content with a maximum of 512 characters.</td></tr><tr><td>COOKIE (cookie)</td><td>Parameters are not supported.</td><td>empty (Content is empty.)<br>null (do not exist)<br>rematch (regular expression matching)</td><td><font color="Red">Unsupported currently</font></td></tr><tr><td>GET_COOKIES_NAMES (cookie parameter name)</td><td>Parameters are not supported.</td><td>Same logical symbol as the matching field <font color="Red">GET Parameter Name</font></td><td>Enter the content with a maximum of 512 characters.</td></tr><tr><td>ARGS_COOKIE (cookie parameter value)</td><td>Parameter entry is supported.</td><td>Same logical symbol as the matching field <font color="Red">GET Parameter Value</font></td><td>Enter the contentwith a maximum of 512 characters.</td></tr><tr><td>GET_HEADERS_NAMES (header parameter name)</td><td>Parameters are not supported.</td><td>exist (Parameter exists.)<br>nexist (Parameter does not exist.)<br>len_eq (length equals to)<br>len_gt (length is greater than)<br>len_lt (length is less than)<br>strprefix (prefix matching)<br>strsuffix (suffix matching)<br>rematch (regular expression matching)</td><td>Enter the content with a maximum of 512 characters. The lowercase is recommended.</td></tr><tr><td>ARGS_HEADER (header parameter value)</td><td>Parameter entry is supported.</td><td>contains (contain)<br>ncontains (do not contain)<br>len_eq (length equals to)<br>len_gt (length is greater than)<br>len_lt (length is less than)<br>strprefix (prefix matching)<br>strsuffix (suffix matching)<br>rematch (regular expression matching)</td><td>Enter the content with a maximum of 512 characters.</td></tr><tr><td>CAPTCHA_RISK (CAPTCHA risk)</td><td>Parameters are not supported.</td><td>eq (equal to)<br>neq (not equal to)<br>belong (belongs to)<br>not_belong (does not belong to)<br>null (does not exist)<br>exist (exists)</td><td>Enter risk level value, supporting numerical range 0-255</td></tr><tr><td>CAPTCHA_DEVICE_RISK (CAPTCHA device risk)</td><td>Parameters are not supported.</td><td>eq (equal to)<br>neq (not equal to)<br>belong (belongs to)<br>not_belong (does not belong to)<br>null (does not exist)<br>exist (exists)</td><td>Enter device risk code, supporting values: 101, 201, 301, 401, 501, 601, 701</td></tr><tr><td>CAPTCHAR_SCORE (CAPTCHA risk assessment score)</td><td>Parameters are not supported.</td><td>numeq (numerically equal to)<br>numgt (numerically greater than)<br>numlt (numerically less than)<br>numle (numerically less than or equal to)<br>numge (numerically greater than or equal to)<br>null (does not exist)<br>exist (exists)</td><td>Enter assessment score, supporting numerical range 0-100</td></tr></tbody></table>
          	  				Note: This field may return null, indicating that no valid values can be obtained.`,
						},
						"compare_func": {
							Type:     schema.TypeString,
							Required: true,
							Description: `Logic symbol
                            Logical symbols are divided into the following types:
								empty (content is empty)
								null (do not exist)
								eq (equal to)
								neq (not equal to)
								contains (contain)
								ncontains (do not contain)
								strprefix (prefix matching)
								strsuffix (suffix matching)
								len_eq (length equals to)
								len_gt (length is greater than)
								len_lt (length is less than)
								ipmatch (belong to)
								ipnmatch (do not belong to)
								numgt (number greater than)
								numlt (number less than)
								geo_in (IP geo belongs to)
								geo_not_in (IP geo not belongs to)
								rematch (regex match)
								numgt (numerically greater than)
								numlt (numerically less than)
								numeq (numerically equal to)
								numneq (numerically not equal to)
								numle (numerically less than or equal to)
								numge (numerically greater than or equal to)
                            Different matching fields correspond to different logical operators. For details, see the matching field table above.
                        Note: This field may return null, indicating that no valid values can be obtained.`,
						},
						"content": {
							Type:     schema.TypeString,
							Required: true,
							Description: `Matching content
                            Currently, when the matching field is COOKIE (cookie), the matching content is not required. In other scenes, the matching content is required.
                        Note: This field may return null, indicating that no valid values can be obtained.`,
						},
						"arg": {
							Type:     schema.TypeString,
							Required: true,
							Description: `Matching parameter
                            There are two types of configuration parameters: unsupported parameters and supported parameters.
                            The matching parameter can be entered only when the matching field is one of the following four. Otherwise, the parameter is not supported.
                                GET (GET parameter value)		
                                POST (POST parameter value)		
                                ARGS_COOKIE (Cookie parameter value)		
                                ARGS_HEADER (Header parameter value)
                        Note: This field may return null, indicating that no valid values can be obtained.`,
						},
						"case_not_sensitive": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "0: case-sensitive, 1: case-insensitive. Note: This field may return null, indicating that no valid values can be obtained.",
						},
					},
				},
			},
			"domain": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Domain name that needs to add policy.",
			},
			"bypass": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The bypass modules are connected by commas between multiple modules. Supported modules ACL (Custom Rules), OWASP (Rule Engine), Webshell (Malicious File Detection), GeoIP (Geographic Block), BWIP (IP Black and White List), CC, BotRPC (BOT Protection), AntiLeakage (Information Leakage Prevention), API (API Security), AI (AI Engine), ip_outo_deny (IP Block), Applet (Mini Program Traffic Risk Control).",
			},
			"status": {
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CUSTOM_WHITE_RULE_STATUS),
				Default:      CUSTOM_WHITE_RULE_STATUS_1,
				Description:  "The status of the switch, 1 is on, 0 is off, default 1.",
			},
			"job_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Rule execution mode: TimedJob indicates scheduled execution. CronJob indicates periodic execution.",
			},
			"job_date_time": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Rule execution time.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timed": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Time parameters for scheduled execution. Note: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_date_time": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Start timestamp, in seconds. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"end_date_time": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "End timestamp, in seconds. Note: This field may return null, indicating that no valid values can be obtained.",
									},
								},
							},
						},
						"cron": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Time parameters for periodic execution. Note: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"days": {
										Type:        schema.TypeSet,
										Optional:    true,
										Description: "Days in each month for execution. Note: This field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
									"w_days": {
										Type:        schema.TypeSet,
										Optional:    true,
										Description: "Days of each week for execution. Note: This field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
									"start_time": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Start time. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"end_time": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "End time. Note: This field may return null, indicating that no valid values can be obtained.",
									},
								},
							},
						},
						"time_t_zone": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Time zone. Note: This field may return null, indicating that no valid values can be obtained.",
						},
					},
				},
			},
			"logical_op": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Logical operator of configuration mode, and/or.",
			},
			"rule_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "rule ID.",
			},
		},
	}
}

func resourceTencentCloudWafCustomWhiteRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_custom_white_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		request       = waf.NewAddCustomWhiteRuleRequest()
		response      = waf.NewAddCustomWhiteRuleResponse()
		statusRequest = waf.NewModifyCustomWhiteRuleStatusRequest()
		domain        string
		ruleIdStr     string
		status        string
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sort_id"); ok {
		request.SortId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("expire_time"); ok {
		request.ExpireTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("strategies"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			strategy := waf.Strategy{}
			if v, ok := dMap["field"]; ok {
				strategy.Field = helper.String(v.(string))
			}

			if v, ok := dMap["compare_func"]; ok {
				strategy.CompareFunc = helper.String(v.(string))
			}

			if v, ok := dMap["content"]; ok {
				strategy.Content = helper.String(v.(string))
			}

			if v, ok := dMap["arg"]; ok {
				strategy.Arg = helper.String(v.(string))
			}

			if v, ok := dMap["case_not_sensitive"]; ok {
				strategy.CaseNotSensitive = helper.IntUint64(v.(int))
			}

			request.Strategies = append(request.Strategies, &strategy)
		}
	}

	if v, ok := d.GetOk("domain"); ok {
		request.Domain = helper.String(v.(string))
		domain = v.(string)
	}

	if v, ok := d.GetOk("bypass"); ok {
		request.Bypass = helper.String(v.(string))
	}

	if v, ok := d.GetOk("job_type"); ok {
		request.JobType = helper.String(v.(string))
	}

	if jobDateTimeMap, ok := helper.InterfacesHeadMap(d, "job_date_time"); ok {
		jobDateTime := waf.JobDateTime{}
		if v, ok := jobDateTimeMap["timed"]; ok {
			for _, item := range v.([]interface{}) {
				timedMap := item.(map[string]interface{})
				timedJob := waf.TimedJob{}
				if v, ok := timedMap["start_date_time"].(int); ok {
					timedJob.StartDateTime = helper.IntUint64(v)
				}

				if v, ok := timedMap["end_date_time"].(int); ok {
					timedJob.EndDateTime = helper.IntUint64(v)
				}

				jobDateTime.Timed = append(jobDateTime.Timed, &timedJob)
			}
		}

		if v, ok := jobDateTimeMap["cron"]; ok {
			for _, item := range v.([]interface{}) {
				cronMap := item.(map[string]interface{})
				cronJob := wafv20180125.CronJob{}
				if v, ok := cronMap["days"]; ok {
					daysSet := v.(*schema.Set).List()
					for i := range daysSet {
						days := daysSet[i].(int)
						cronJob.Days = append(cronJob.Days, helper.IntUint64(days))
					}
				}

				if v, ok := cronMap["w_days"]; ok {
					wDaysSet := v.(*schema.Set).List()
					for i := range wDaysSet {
						wDays := wDaysSet[i].(int)
						cronJob.WDays = append(cronJob.WDays, helper.IntUint64(wDays))
					}
				}

				if v, ok := cronMap["start_time"].(string); ok && v != "" {
					cronJob.StartTime = helper.String(v)
				}

				if v, ok := cronMap["end_time"].(string); ok && v != "" {
					cronJob.EndTime = helper.String(v)
				}

				jobDateTime.Cron = append(jobDateTime.Cron, &cronJob)
			}
		}

		if v, ok := jobDateTimeMap["time_t_zone"].(string); ok && v != "" {
			jobDateTime.TimeTZone = helper.String(v)
		}

		request.JobDateTime = &jobDateTime
	}

	if v, ok := d.GetOk("logical_op"); ok {
		request.LogicalOp = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().AddCustomWhiteRule(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create waf CustomWhiteRule failed, reason:%+v", logId, err)
		return err
	}

	ruleId := *response.Response.RuleId
	ruleIdStr = strconv.FormatUint(ruleId, 10)

	if v, ok := d.GetOk("status"); ok {
		status = v.(string)
	}

	if status == CUSTOM_WHITE_RULE_STATUS_0 {
		statusRequest.Domain = &domain
		statusRequest.RuleId = &ruleId
		statusRequest.Status = helper.IntUint64(CUSTOM_WHITE_RULE_STATUS_0_INT)
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().ModifyCustomWhiteRuleStatus(statusRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s modify waf CustomWhiteRule status failed, reason:%+v", logId, err)
			return err
		}
	}

	d.SetId(strings.Join([]string{domain, ruleIdStr}, tccommon.FILED_SP))
	return resourceTencentCloudWafCustomWhiteRuleRead(d, meta)
}

func resourceTencentCloudWafCustomWhiteRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_custom_white_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	domain := idSplit[0]
	ruleId := idSplit[1]

	customWhiteRule, err := service.DescribeWafCustomWhiteRuleById(ctx, domain, ruleId)
	if err != nil {
		return err
	}

	if customWhiteRule == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `WafCustomWhiteRule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if customWhiteRule.Name != nil {
		_ = d.Set("name", customWhiteRule.Name)
	}

	if customWhiteRule.SortId != nil {
		_ = d.Set("sort_id", customWhiteRule.SortId)
	}

	if customWhiteRule.ExpireTime != nil {
		_ = d.Set("expire_time", customWhiteRule.ExpireTime)
	}

	if customWhiteRule.Strategies != nil {
		strategiesList := []interface{}{}
		for _, strategies := range customWhiteRule.Strategies {
			strategiesMap := map[string]interface{}{}

			if strategies.Field != nil {
				strategiesMap["field"] = strategies.Field
			}

			if strategies.CompareFunc != nil {
				strategiesMap["compare_func"] = strategies.CompareFunc
			}

			if strategies.Content != nil {
				strategiesMap["content"] = strategies.Content
			}

			if strategies.Arg != nil {
				strategiesMap["arg"] = strategies.Arg
			}

			if strategies.CaseNotSensitive != nil {
				strategiesMap["case_not_sensitive"] = strategies.CaseNotSensitive
			}

			strategiesList = append(strategiesList, strategiesMap)
		}

		_ = d.Set("strategies", strategiesList)

	}

	_ = d.Set("domain", domain)

	if customWhiteRule.JobType != nil {
		_ = d.Set("job_type", customWhiteRule.JobType)
	}

	if customWhiteRule.JobDateTime != nil {
		tmpList := make([]map[string]interface{}, 0)
		dMap := map[string]interface{}{}
		if customWhiteRule.JobDateTime.Timed != nil {
			timedList := []interface{}{}
			for _, v := range customWhiteRule.JobDateTime.Timed {
				timedMap := map[string]interface{}{}
				if v.StartDateTime != nil {
					timedMap["start_date_time"] = v.StartDateTime
				}

				if v.EndDateTime != nil {
					timedMap["end_date_time"] = v.EndDateTime
				}

				timedList = append(timedList, timedMap)
			}

			dMap["timed"] = timedList
		}

		if customWhiteRule.JobDateTime.Cron != nil {
			cronList := []interface{}{}
			for _, v := range customWhiteRule.JobDateTime.Cron {
				cronMap := map[string]interface{}{}
				if v.Days != nil {
					cronMap["days"] = v.Days
				}

				if v.WDays != nil {
					cronMap["w_days"] = v.WDays
				}

				if v.StartTime != nil {
					cronMap["start_time"] = v.StartTime
				}

				if v.EndTime != nil {
					cronMap["end_time"] = v.EndTime
				}

				cronList = append(cronList, cronMap)
			}

			dMap["cron"] = cronList
		}

		if customWhiteRule.JobDateTime.TimeTZone != nil {
			dMap["time_t_zone"] = customWhiteRule.JobDateTime.TimeTZone
		}

		tmpList = append(tmpList, dMap)

		_ = d.Set("job_date_time", tmpList)
	}

	if customWhiteRule.Bypass != nil {
		_ = d.Set("bypass", customWhiteRule.Bypass)
	}

	if customWhiteRule.Status != nil {
		_ = d.Set("status", customWhiteRule.Status)
	}

	if customWhiteRule.LogicalOp != nil {
		_ = d.Set("logical_op", customWhiteRule.LogicalOp)
	}

	if customWhiteRule.RuleId != nil {
		_ = d.Set("rule_id", customWhiteRule.RuleId)
	}

	return nil
}

func resourceTencentCloudWafCustomWhiteRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_custom_white_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		request       = waf.NewModifyCustomWhiteRuleRequest()
		statusRequest = waf.NewModifyCustomWhiteRuleStatusRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	domain := idSplit[0]
	ruleId := idSplit[1]

	immutableArgs := []string{"domain"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	request.Domain = &domain
	ruleIdInt, _ := strconv.ParseInt(ruleId, 10, 64)
	ruleIdUInt := uint64(ruleIdInt)
	request.RuleId = &ruleIdUInt

	if v, ok := d.GetOk("name"); ok {
		request.RuleName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("bypass"); ok {
		request.Bypass = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sort_id"); ok {
		tmpSortId, _ := strconv.ParseInt(v.(string), 10, 64)
		request.SortId = helper.Int64Uint64(tmpSortId)
	}

	if v, ok := d.GetOk("expire_time"); ok {
		tmpExpireTime, _ := strconv.ParseInt(v.(string), 10, 64)
		request.ExpireTime = helper.Int64Uint64(tmpExpireTime)
	}

	if v, ok := d.GetOk("strategies"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			strategy := waf.Strategy{}
			if v, ok := dMap["field"]; ok {
				strategy.Field = helper.String(v.(string))
			}

			if v, ok := dMap["compare_func"]; ok {
				strategy.CompareFunc = helper.String(v.(string))
			}

			if v, ok := dMap["content"]; ok {
				strategy.Content = helper.String(v.(string))
			}

			if v, ok := dMap["arg"]; ok {
				strategy.Arg = helper.String(v.(string))
			}

			if v, ok := dMap["case_not_sensitive"]; ok {
				strategy.CaseNotSensitive = helper.IntUint64(v.(int))
			}

			request.Strategies = append(request.Strategies, &strategy)
		}
	}

	if v, ok := d.GetOk("job_type"); ok {
		request.JobType = helper.String(v.(string))
	}

	if jobDateTimeMap, ok := helper.InterfacesHeadMap(d, "job_date_time"); ok {
		jobDateTime := waf.JobDateTime{}
		if v, ok := jobDateTimeMap["timed"]; ok {
			for _, item := range v.([]interface{}) {
				timedMap := item.(map[string]interface{})
				timedJob := waf.TimedJob{}
				if v, ok := timedMap["start_date_time"].(int); ok {
					timedJob.StartDateTime = helper.IntUint64(v)
				}

				if v, ok := timedMap["end_date_time"].(int); ok {
					timedJob.EndDateTime = helper.IntUint64(v)
				}

				jobDateTime.Timed = append(jobDateTime.Timed, &timedJob)
			}
		}

		if v, ok := jobDateTimeMap["cron"]; ok {
			for _, item := range v.([]interface{}) {
				cronMap := item.(map[string]interface{})
				cronJob := wafv20180125.CronJob{}
				if v, ok := cronMap["days"]; ok {
					daysSet := v.(*schema.Set).List()
					for i := range daysSet {
						days := daysSet[i].(int)
						cronJob.Days = append(cronJob.Days, helper.IntUint64(days))
					}
				}

				if v, ok := cronMap["w_days"]; ok {
					wDaysSet := v.(*schema.Set).List()
					for i := range wDaysSet {
						wDays := wDaysSet[i].(int)
						cronJob.WDays = append(cronJob.WDays, helper.IntUint64(wDays))
					}
				}

				if v, ok := cronMap["start_time"].(string); ok && v != "" {
					cronJob.StartTime = helper.String(v)
				}

				if v, ok := cronMap["end_time"].(string); ok && v != "" {
					cronJob.EndTime = helper.String(v)
				}

				jobDateTime.Cron = append(jobDateTime.Cron, &cronJob)
			}
		}

		if v, ok := jobDateTimeMap["time_t_zone"].(string); ok && v != "" {
			jobDateTime.TimeTZone = helper.String(v)
		}

		request.JobDateTime = &jobDateTime
	}

	if v, ok := d.GetOk("logical_op"); ok {
		request.LogicalOp = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().ModifyCustomWhiteRule(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update waf CustomWhiteRule failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("status") {
		if v, ok := d.GetOk("status"); ok {
			status := v.(string)
			statusRequest.Domain = &domain
			statusRequest.RuleId = &ruleIdUInt
			if status == CUSTOM_WHITE_RULE_STATUS_0 {
				statusRequest.Status = helper.IntUint64(CUSTOM_WHITE_RULE_STATUS_0_INT)
			} else {
				statusRequest.Status = helper.IntUint64(CUSTOM_WHITE_RULE_STATUS_1_INT)
			}

			err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().ModifyCustomWhiteRuleStatus(statusRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s modify waf CustomRule status failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	return resourceTencentCloudWafCustomWhiteRuleRead(d, meta)
}

func resourceTencentCloudWafCustomWhiteRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_custom_white_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	domain := idSplit[0]
	ruleId := idSplit[1]

	if err := service.DeleteWafCustomWhiteRuleById(ctx, domain, ruleId); err != nil {
		return err
	}

	return nil
}
