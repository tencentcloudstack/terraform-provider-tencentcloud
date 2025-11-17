package waf

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wafv20180125 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWafOwaspWhiteRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafOwaspWhiteRuleCreate,
		Read:   resourceTencentCloudWafOwaspWhiteRuleRead,
		Update: resourceTencentCloudWafOwaspWhiteRuleUpdate,
		Delete: resourceTencentCloudWafOwaspWhiteRuleDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule name.",
			},

			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Domain name.",
			},

			"strategies": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Rule-Based matching policy list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the matching field.\n\nDifferent matching fields result in different matching parameters, logical operators, and matching contents. the details are as follows:.\n<table><thead><tr><th>Matching Field</th> <th>Matching Parameter</th> <th>Logical Symbol</th> <th>Matching Content</th></tr></thead> <tbody><tr><td>IP (source IP)</td> <td>Parameters are not supported.</td> <td>ipmatch (match)<br/>ipnmatch (mismatch)</td> <td>Multiple IP addresses are separated by commas. A maximum of 20 IP addresses are allowed.</td></tr> <tr><td>IPv6 (source IPv6)</td> <td>Parameters are not supported.</td> <td>ipmatch (match)<br/>ipnmatch (mismatch)</td> <td>A single IPv6 address is supported.</td></tr> <tr><td>Referer (referer)</td> <td>Parameters are not supported.</td> <td>empty (Content is empty.)<br/>null (do not exist)<br/>eq (equal to)<br/>neq (not equal to)<br/>contains (contain)<br/>ncontains (do not contain)<br/>len_eq (length equals to)<br/>len_gt (length is greater than)<br/>len_lt (length is less than)<br/>strprefix (prefix matching)<br/>strsuffix (suffix matching)<br/>rematch (regular expression matching)</td> <td>Enter the content, with a maximum of 512 characters.</td></tr> <tr><td>URL (request path)</td> <td>Parameters are not supported.</td> <td>eq (equal to)<br/>neq (not equal to)<br/>contains (contain)<br/>ncontains (do not contain)<br/>len_eq (length equals to)<br/>len_gt (length is greater than)<br/>len_lt (length is \n less than)<br/>strprefix (prefix matching)<br/>strsuffix (suffix matching)<br/>rematch (regular expression matching)</td> <td>Enter the content starting with /, with a maximum of 512 characters.</td></tr> <tr><td>UserAgent (UserAgent)</td> <td>Parameters are not supported.</td><td>Same logical symbols as the matching field <font color=\"Red\">Referer</font></td> <td>Enter the content with a maximum of 512 characters.</td></tr> <tr><td>HTTP_METHOD (HTTP request method)</td> <td>Parameters are not supported.</td> <td>eq (equal to)<br/>neq (not equal to)</td> <td>Enter the method name. The uppercase is recommended.</td></tr> <tr><td>QUERY_STRING (request string)</td> <td>Parameters are not supported.</td> <td>Same logical symbol as the matching field <font color=\"Red\">Request Path</font></td><td>Enter the content with a maximum of 512 characters.</td></tr> <tr><td>GET (GET parameter value)</td> <td>Parameter entry is supported.</td> <td>contains (contain)<br/>ncontains (do not contain)<br/>len_eq (length equals to)<br/>len_gt (length is greater than)<br/>len_lt (length is less than)<br/>strprefix (prefix matching)<br/>strsuffix (suffix matching)</td> <td>Enter the content with a maximum of 512 characters.</td></tr> <tr><td>GET_PARAMS_NAMES (GET parameter name)</td> <td>Parameters are not supported.</td> <td>exist (Parameter exists.)<br/>nexist (Parameter does not exist.)<br/>len_eq (length equals to)<br/>len_gt (length is greater than)<br/>len_lt (length is less than)<br/>strprefix (prefix matching)<br/>strsuffix (suffix matching)</td><td>Enter the content with a maximum of 512 characters.</td></tr> <tr><td>POST (POST parameter value)</td> <td>Parameter entry is supported.</td> <td>Same logical symbol as the matching field <font color=\"Red\">GET Parameter Value</font></td> <td>Enter the content with a maximum of 512 characters.</td></tr> <tr><td>GET_POST_NAMES (POST parameter name)</td> <td>Parameters are not supported.</td> <td>Same logical symbol as the matching field <font color=\"Red\">GET Parameter Name</font></td> <td>Enter the content with a maximum of 512 characters.</td></tr> <tr><td>POST_BODY (complete body)</td> <td>Parameters are not supported.</td> <td>Same logical symbol as the matching field <font color=\"Red\">Request Path</font></td><td>Enter the body content with a maximum of 512 characters.</td></tr> <tr><td>COOKIE (cookie)</td> <td>Parameters are not supported.</td> <td>empty (Content is empty.)<br/>null (do not exist)<br/>rematch (regular expression matching)</td> <td><font color=\"Red\">Unsupported currently</font></td></tr> <tr><td>GET_COOKIES_NAMES (cookie parameter name)</td> <td>Parameters are not supported.</td> <td>Same logical symbol as the matching field <font color=\"Red\">GET Parameter Name</font></td> <td>Enter the content with a maximum of 512 characters.</td></tr> <tr><td>ARGS_COOKIE (cookie parameter value)</td> <td>Parameter entry is supported.</td> <td>Same logical symbol as the matching field <font color=\"Red\">GET Parameter Value</font></td> <td>Enter the content512 characters limit</td></tr><tr><td>GET_HEADERS_NAMES (Header parameter name)</td><td>parameter not supported</td><td>exsit (parameter exists)<br/>nexsit (parameter does not exist)<br/>len_eq (LENGTH equal)<br/>len_gt (LENGTH greater than)<br/>len_lt (LENGTH less than)<br/>strprefix (prefix match)<br/>strsuffix (suffix matching)<br/>rematch (regular expression matching)</td><td>enter CONTENT, lowercase is recommended, up to 512 characters</td></tr><tr><td>ARGS_Header (Header parameter value)</td><td>support parameter entry</td><td>contains (include)<br/>ncontains (does not include)<br/>len_eq (LENGTH equal)<br/>len_gt (LENGTH greater than)<br/>len_lt (LENGTH less than)<br/>strprefix (prefix match)<br/>strsuffix (suffix matching)<br/>rematch (regular expression matching)</td><td>enter CONTENT, up to 512 characters</td></tr><tr><td>CONTENT_LENGTH (CONTENT-LENGTH)</td><td>support parameter entry</td><td>numgt (value greater than)<br/>numlt (value smaller than)<br/>numeq (value equal to)<br/></td><td>enter an integer between 0-9999999999999</td></tr><tr><td>IP_GEO (source IP geolocation)</td><td>support parameter entry</td><td>GEO_in (belong)<br/>GEO_not_in (not_in)<br/></td><td>enter CONTENT, up to 10240 characters, format: serialized JSON, format: [{\"Country\":\"china\",\"Region\":\"guangdong\",\"City\":\"shenzhen\"}]</td></tr><tr><td>CAPTCHA_RISK (CAPTCHA RISK)</td><td>parameter not supported</td><td>eq (equal)<br/>neq (not equal to)<br/>belong (belong)<br/>not_belong (not belong to)<br/>null (nonexistent)<br/>exist (exist)</td><td>enter RISK level value, value range 0-255</td></tr><tr><td>CAPTCHA_DEVICE_RISK (CAPTCHA DEVICE RISK)</td><td>parameter not supported</td><td>eq (equal)<br/>neq (not equal to)<br/>belong (belong)<br/>not_belong (not belong to)<br/>null (nonexistent)<br/>exist (exist)</td><td>enter DEVICE RISK code, valid values: 101, 201, 301, 401, 501, 601, 701</td></tr><tr><td>CAPTCHAR_SCORE (CAPTCHA RISK assessment SCORE)</td><td>parameter not supported</td><td>numeq (value equal to)<br/>numgt (value greater than)<br/>numlt (value smaller than)<br/>numle (less than or equal to)<br/>numge (value is greater than or equal to)<br/>null (nonexistent)<br/>exist (exist)</td><td>enter assessment SCORE, value range 0-100</td></tr>.\n</tbody></table>.",
						},
						"compare_func": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the logic symbol. \n\nLogical symbols are divided into the following types:.\nEmpty (content is empty).\nnull (not found).\nEq (equal to).\nneq (not equal to).\ncontains (contain).\nncontains (do not contain).\nstrprefix (prefix matching).\nstrsuffix (suffix matching).\nLen_eq (length equals to).\nLen_gt (length greater than).\nLen_lt (length less than).\nipmatch (belong).\nipnmatch (not_in).\nnumgt (value greater than).\nNumValue smaller than].\nValue equal to.\nnumneq (value not equal to).\nnumle (less than or equal to).\nnumge (value is greater than or equal to).\ngeo_in (IP geographic belong).\ngeo_not_in (IP geographic not_in).\nSpecifies different logical operators for matching fields. for details, see the matching field table above.",
						},
						"content": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the match content.\n\nCurrently, when the match field is COOKIE (COOKIE), match content is not required. all others are needed.",
						},
						"arg": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the matching parameter.\n\nConfiguration parameters are divided into two data types: parameter not supported and support parameters.\nWhen the match field is one of the following four, the matching parameter can be entered, otherwise not supported.\nGET (get parameter value).\v\t\t\nPOST (post parameter value).\v\t\t\nARGS_COOKIE (COOKIE parameter value).\v\t\t\nARGS_HEADER (HEADER parameter value).",
						},
						"case_not_sensitive": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Case-Sensitive.\nCase-Insensitive.",
						},
					},
				},
			},

			"ids": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "ID list of allowlisted rules.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Allowlist type. valid values: 0 (allowlisting by specific rule ID), 1 (allowlisting by rule type).",
			},

			"job_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule execution mode: TimedJob indicates scheduled execution. CronJob indicates periodic execution.",
			},

			"job_date_time": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Scheduled task configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timed": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Time parameter for scheduled execution.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_date_time": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Start timestamp, in seconds.",
									},
									"end_date_time": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "End timestamp, in seconds.",
									},
								},
							},
						},
						"cron": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Time parameter for periodic execution.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"days": {
										Type:        schema.TypeSet,
										Optional:    true,
										Description: "Execution day of each month.",
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
									"w_days": {
										Type:        schema.TypeSet,
										Optional:    true,
										Description: "Execution day of each week.",
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
									"start_time": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Start time.",
									},
									"end_time": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "End time.",
									},
								},
							},
						},
						"time_t_zone": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the time zone.",
						},
					},
				},
			},

			"expire_time": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "If the JobDateTime field is not set, this field is used. 0 means permanent, other values indicate the cutoff time for scheduled effect (unit: seconds).",
			},

			"status": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Rule status. valid values: 0 (disabled), 1 (enabled). enabled by default.",
			},

			// computed
			"rule_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Rule ID.",
			},
		},
	}
}

func resourceTencentCloudWafOwaspWhiteRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_owasp_white_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = wafv20180125.NewCreateOwaspWhiteRuleRequest()
		response = wafv20180125.NewCreateOwaspWhiteRuleResponse()
		domain   string
		ruleId   string
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("domain"); ok {
		request.Domain = helper.String(v.(string))
		domain = v.(string)
	}

	if v, ok := d.GetOk("strategies"); ok {
		for _, item := range v.([]interface{}) {
			strategiesMap := item.(map[string]interface{})
			strategy := wafv20180125.Strategy{}
			if v, ok := strategiesMap["field"].(string); ok && v != "" {
				strategy.Field = helper.String(v)
			}

			if v, ok := strategiesMap["compare_func"].(string); ok && v != "" {
				strategy.CompareFunc = helper.String(v)
			}

			if v, ok := strategiesMap["content"].(string); ok && v != "" {
				strategy.Content = helper.String(v)
			}

			if v, ok := strategiesMap["arg"].(string); ok && v != "" {
				strategy.Arg = helper.String(v)
			}

			if v, ok := strategiesMap["case_not_sensitive"].(int); ok {
				strategy.CaseNotSensitive = helper.IntUint64(v)
			}

			request.Strategies = append(request.Strategies, &strategy)
		}
	}

	if v, ok := d.GetOk("ids"); ok {
		idsSet := v.(*schema.Set).List()
		for i := range idsSet {
			ids := idsSet[i].(int)
			request.Ids = append(request.Ids, helper.IntUint64(ids))
		}
	}

	if v, ok := d.GetOkExists("type"); ok {
		request.Type = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("job_type"); ok {
		request.JobType = helper.String(v.(string))
	}

	if jobDateTimeMap, ok := helper.InterfacesHeadMap(d, "job_date_time"); ok {
		jobDateTime := wafv20180125.JobDateTime{}
		if v, ok := jobDateTimeMap["timed"]; ok {
			for _, item := range v.([]interface{}) {
				timedMap := item.(map[string]interface{})
				timedJob := wafv20180125.TimedJob{}
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

	if v, ok := d.GetOkExists("expire_time"); ok {
		request.ExpireTime = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("status"); ok {
		request.Status = helper.IntUint64(v.(int))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().CreateOwaspWhiteRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create waf owasp white rule failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create waf owasp white rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.RuleId == nil {
		return fmt.Errorf("RuleId is nil.")
	}

	ruleId = helper.UInt64ToStr(*response.Response.RuleId)
	d.SetId(strings.Join([]string{domain, ruleId}, tccommon.FILED_SP))
	return resourceTencentCloudWafOwaspWhiteRuleRead(d, meta)
}

func resourceTencentCloudWafOwaspWhiteRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_owasp_white_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	domain := idSplit[0]
	ruleId := idSplit[1]

	respData, err := service.DescribeWafOwaspWhiteRuleById(ctx, domain, ruleId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_waf_owasp_white_rule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("domain", domain)

	if respData.Name != nil {
		_ = d.Set("name", respData.Name)
	}

	if respData.Strategies != nil {
		tmpList := make([]map[string]interface{}, 0, len(respData.Strategies))
		for _, strategies := range respData.Strategies {
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

			tmpList = append(tmpList, strategiesMap)
		}

		_ = d.Set("strategies", tmpList)
	}

	if respData.Ids != nil {
		_ = d.Set("ids", respData.Ids)
	}

	if respData.Type != nil {
		_ = d.Set("type", respData.Type)
	}

	if respData.JobType != nil {
		_ = d.Set("job_type", respData.JobType)
	}

	if respData.JobDateTime != nil {
		jobDateTimeMap := map[string]interface{}{}
		if respData.JobDateTime.Timed != nil {
			timedList := make([]map[string]interface{}, 0, len(respData.JobDateTime.Timed))
			for _, timed := range respData.JobDateTime.Timed {
				timedMap := map[string]interface{}{}
				if timed.StartDateTime != nil {
					timedMap["start_date_time"] = timed.StartDateTime
				}

				if timed.EndDateTime != nil {
					timedMap["end_date_time"] = timed.EndDateTime
				}

				timedList = append(timedList, timedMap)
			}

			jobDateTimeMap["timed"] = timedList
		}

		if respData.JobDateTime.Cron != nil {
			cronList := make([]map[string]interface{}, 0, len(respData.JobDateTime.Cron))
			for _, cron := range respData.JobDateTime.Cron {
				cronMap := map[string]interface{}{}
				if cron.Days != nil {
					cronMap["days"] = cron.Days
				}

				if cron.WDays != nil {
					cronMap["w_days"] = cron.WDays
				}

				if cron.StartTime != nil {
					cronMap["start_time"] = cron.StartTime
				}

				if cron.EndTime != nil {
					cronMap["end_time"] = cron.EndTime
				}

				cronList = append(cronList, cronMap)
			}

			jobDateTimeMap["cron"] = cronList
		}

		if respData.JobDateTime.TimeTZone != nil {
			jobDateTimeMap["time_t_zone"] = respData.JobDateTime.TimeTZone
		}

		_ = d.Set("job_date_time", []interface{}{jobDateTimeMap})
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	if respData.RuleId != nil {
		_ = d.Set("rule_id", respData.RuleId)
	}

	return nil
}

func resourceTencentCloudWafOwaspWhiteRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_owasp_white_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	domain := idSplit[0]
	ruleId := idSplit[1]

	needChange := false
	mutableArgs := []string{"name", "strategies", "ids", "type", "job_type", "job_date_time", "expire_time", "status"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := wafv20180125.NewModifyOwaspWhiteRuleRequest()
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOk("strategies"); ok {
			for _, item := range v.([]interface{}) {
				strategiesMap := item.(map[string]interface{})
				strategy := wafv20180125.Strategy{}
				if v, ok := strategiesMap["field"].(string); ok && v != "" {
					strategy.Field = helper.String(v)
				}

				if v, ok := strategiesMap["compare_func"].(string); ok && v != "" {
					strategy.CompareFunc = helper.String(v)
				}

				if v, ok := strategiesMap["content"].(string); ok && v != "" {
					strategy.Content = helper.String(v)
				}

				if v, ok := strategiesMap["arg"].(string); ok && v != "" {
					strategy.Arg = helper.String(v)
				}

				if v, ok := strategiesMap["case_not_sensitive"].(int); ok {
					strategy.CaseNotSensitive = helper.IntUint64(v)
				}

				request.Strategies = append(request.Strategies, &strategy)
			}
		}

		if v, ok := d.GetOk("ids"); ok {
			idsSet := v.(*schema.Set).List()
			for i := range idsSet {
				ids := idsSet[i].(int)
				request.Ids = append(request.Ids, helper.IntUint64(ids))
			}
		}

		if v, ok := d.GetOkExists("type"); ok {
			request.Type = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("job_type"); ok {
			request.JobType = helper.String(v.(string))
		}

		if jobDateTimeMap, ok := helper.InterfacesHeadMap(d, "job_date_time"); ok {
			jobDateTime := wafv20180125.JobDateTime{}
			if v, ok := jobDateTimeMap["timed"]; ok {
				for _, item := range v.([]interface{}) {
					timedMap := item.(map[string]interface{})
					timedJob := wafv20180125.TimedJob{}
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

		if v, ok := d.GetOkExists("expire_time"); ok {
			request.ExpireTime = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("status"); ok {
			request.Status = helper.IntUint64(v.(int))
		}

		request.Domain = &domain
		request.RuleId = helper.StrToUint64Point(ruleId)
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().ModifyOwaspWhiteRuleWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update waf owasp white rule failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudWafOwaspWhiteRuleRead(d, meta)
}

func resourceTencentCloudWafOwaspWhiteRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_owasp_white_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = wafv20180125.NewDeleteOwaspWhiteRuleRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	domain := idSplit[0]
	ruleId := idSplit[1]

	request.Domain = &domain
	request.Ids = append(request.Ids, helper.StrToUint64Point(ruleId))
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().DeleteOwaspWhiteRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete waf owasp white rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
