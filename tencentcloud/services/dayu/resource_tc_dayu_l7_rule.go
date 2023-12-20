package dayu

import (
	"context"
	"fmt"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dayu "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dayu/v20180709"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDayuL7Rule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDayuL7RuleCreate,
		Read:   resourceTencentCloudDayuL7RuleRead,
		Update: resourceTencentCloudDayuL7RuleUpdate,
		Delete: resourceTencentCloudDayuL7RuleDelete,

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the resource that the layer 7 rule works for.",
			},
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_RESOURCE_TYPE),
				ForceNew:     true,
				Description:  "Type of the resource that the layer 7 rule works for, valid value is `bgpip`.",
			},
			"domain": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(0, 80),
				Description:  "Domain that the layer 7 rule works for. Valid string length ranges from 0 to 80.",
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_L7_RULE_PROTOCOL),
				Description:  "Protocol of the rule. Valid values: `http`, `https`.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the rule.",
			},
			"switch": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicate the rule will take effect or not.",
			},
			"source_type": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue(DAYU_L7_RULE_SOURCE_TYPE),
				Description:  "Source type, `1` for source of host, `2` for source of IP.",
			},
			"source_list": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type:        schema.TypeString,
					Description: "Source IP or domain, valid format of ip is like `1.1.1.1:60` or `1.1.1.1` and valid format of host source is like `abc.com` or `abc.com:80`.",
				},
				MinItems:    1,
				MaxItems:    16,
				Description: "Source list of the rule, it can be a set of ip sources or a set of domain sources. The number of items ranges from 1 to 16.",
			},
			"ssl_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SSL ID, when the `protocol` is `https`, the field should be set with valid SSL id.",
			},
			"health_check_switch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates whether health check is enabled. The default is `false`.",
			},
			"health_check_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(10, 60),
				Description:  "Interval time of health check. Valid value ranges: [10~60]sec. The default is 15 sec.",
			},
			"health_check_health_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(2, 10),
				Description:  "Health threshold of health check, and the default is `3`. If a success result is returned for the health check 3 consecutive times, indicates that the forwarding is normal. The value range is [2-10].",
			},
			"health_check_unhealth_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(2, 10),
				Description:  "Unhealthy threshold of health check, and the default is `3`. If the unhealthy result is returned 3 consecutive times, indicates that the forwarding is abnormal. The value range is [2-10].",
			},
			"health_check_code": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(1, 31),
				Description:  "HTTP Status Code. The default is `26`. Valid value ranges: [1~31]. `1` means the return value '1xx' is health. `2` means the return value '2xx' is health. `4` means the return value '3xx' is health. `8` means the return value '4xx' is health. `16` means the return value '5xx' is health. If you want multiple return codes to indicate health, need to add the corresponding values.",
			},
			"health_check_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Path of health check. The default is `/`.",
			},
			"health_check_method": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_RULE_METHOD),
				Description:  "Methods of health check. The default is 'HEAD', the available value are 'HEAD' and 'GET'.",
			},
			//computed
			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the layer 7 rule.",
			},
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Status of the rule. `0` for create/modify success, `2` for create/modify fail, `3` for delete success, `5` for delete failed, `6` for waiting to be created/modified, `7` for waiting to be deleted and 8 for waiting to get SSL ID.",
			},
		},
	}
}

func resourceTencentCloudDayuL7RuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_l7_rule.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	resourceId := d.Get("resource_id").(string)
	resourceType := d.Get("resource_type").(string)

	domain := d.Get("domain").(string)

	//set L4RuleEntry
	var rule dayu.L7RuleEntry
	rule.LbType = helper.IntUint64(1)
	//test that the keep time para will make effect
	protocol := d.Get("protocol").(string)
	sslId := ""
	if protocol == "https" {
		if v, ok := d.GetOk("ssl_id"); ok {
			sslId = v.(string)
		}
		if sslId == "" {
			return fmt.Errorf("`ssl_id` should be set when `protocol` is `https`.")
		}
		rule.SSLId = &sslId
		rule.CertType = helper.IntUint64(2)
	} else {
		rule.CertType = helper.IntUint64(0)
	}
	rule.Protocol = &protocol
	rule.RuleName = helper.String(d.Get("name").(string))
	sourceType := d.Get("source_type").(int)
	//check that there is no check with the source list and sdk returns
	rule.SourceType = helper.IntUint64(sourceType)
	rule.Domain = &domain

	sourceList := d.Get("source_list").(*schema.Set).List()
	//check
	healthCheckSwitch := d.Get("health_check_switch").(bool)
	if healthCheckSwitch {
		if len(sourceList) <= 1 {
			return fmt.Errorf("The `health_check_switch` cannot be set `true` when `source_list` has only one item.")
		}
	}

	for _, source := range sourceList {
		var l4RuleSource dayu.L4RuleSource
		l4RuleSource.Source = helper.String(source.(string))
		l4RuleSource.Weight = helper.IntUint64(0)
		rule.SourceList = append(rule.SourceList, &l4RuleSource)
	}

	dayuService := DayuService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	ruleId := ""
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := dayuService.CreateL7Rule(ctx, resourceType, resourceId, rule)
		if e != nil {
			return tccommon.RetryError(e)
		}
		ruleId = result
		return nil
	})

	if err != nil {
		return err
	}

	d.SetId(resourceType + tccommon.FILED_SP + resourceId + tccommon.FILED_SP + ruleId)

	readyFlag, rErr := checkL7RuleStatus(meta, resourceType, resourceId, ruleId, "create")
	if rErr != nil {
		return rErr
	}
	if !readyFlag {
		return fmt.Errorf("Creating operation is timeout %s", ruleId)
	}

	//set health check
	var healthCheck dayu.L7HealthConfig
	healthCheck.Protocol = helper.String(d.Get("protocol").(string))
	healthCheck.Domain = &domain
	healthCheck.Enable = helper.BoolToInt64Pointer(d.Get("health_check_switch").(bool))
	healthCheck.Interval = helper.IntUint64(d.Get("health_check_interval").(int))
	healthCheck.Method = helper.String(d.Get("health_check_method").(string))
	healthCheck.Url = helper.String(d.Get("health_check_path").(string))
	healthCheck.KickNum = helper.IntUint64(d.Get("health_check_unhealth_num").(int))
	healthCheck.AliveNum = helper.IntUint64(d.Get("health_check_health_num").(int))
	healthCheck.StatusCode = helper.IntUint64(d.Get("health_check_code").(int))

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := dayuService.SetL7Health(ctx, resourceType, resourceId, healthCheck)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})

	if err != nil {
		return err
	}

	readyFlag, rErr = checkL7RuleStatus(meta, resourceType, resourceId, ruleId, "check_health")
	if rErr != nil {
		return rErr
	}
	if !readyFlag {
		return fmt.Errorf("Set health is timeout %s", ruleId)
	}

	//set switch
	switchFlag := d.Get("switch").(bool)

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := dayuService.SetRuleSwitch(ctx, resourceType, resourceId, ruleId, switchFlag, protocol)
		if e != nil {
			return tccommon.RetryError(e, tccommon.InternalError)
		}
		return nil
	})

	if err != nil {
		return err
	}

	//check switch status
	readyFlag, rErr = checkL7RuleStatus(meta, resourceType, resourceId, ruleId, fmt.Sprintf("check_switch_%t", switchFlag))
	if rErr != nil {
		return rErr
	}
	if !readyFlag {
		return fmt.Errorf("Set switch is timeout %s", ruleId)
	}

	return resourceTencentCloudDayuL7RuleRead(d, meta)
}

func resourceTencentCloudDayuL7RuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_l7_rule.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 3 {
		return fmt.Errorf("broken ID of dayu L7 rule")
	}
	resourceType := items[0]
	resourceId := items[1]
	ruleId := items[2]

	domain := d.Get("domain").(string)
	sourceList := d.Get("source_list").(*schema.Set).List()

	d.Partial(true)
	dayuService := DayuService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	ruleFlag := false
	ruleKey := []string{"protocol", "source_type", "source_list", "ssl_id"}

	for _, key := range ruleKey {
		if d.HasChange(key) {
			ruleFlag = true
		}
	}
	if ruleFlag {
		//set L4RuleEntry
		var rule dayu.L7RuleEntry
		rule.LbType = helper.IntUint64(1)
		rule.RuleId = helper.String(ruleId)
		//test that the keep time para will make effect
		protocol := d.Get("protocol").(string)
		sslId := ""
		if protocol == DAYU_L7_RULE_PROTOCOL_HTTPS {
			if v, ok := d.GetOk("ssl_id"); ok {
				sslId = v.(string)
			}
			if sslId == "" {
				return fmt.Errorf("`ssl_id` should be set when `protocol` is `https`.")
			}
			rule.SSLId = &sslId
			rule.CertType = helper.IntUint64(2)
		} else {
			rule.CertType = helper.IntUint64(0)
		}
		rule.RuleName = helper.String(d.Get("name").(string))
		sourceType := d.Get("source_type").(int)
		//check that there is no check with the source list and sdk returns
		rule.SourceType = helper.IntUint64(sourceType)
		rule.Domain = &domain
		rule.Protocol = &protocol

		for _, source := range sourceList {
			var l4RuleSource dayu.L4RuleSource
			l4RuleSource.Source = helper.String(source.(string))
			l4RuleSource.Weight = helper.IntUint64(0)
			rule.SourceList = append(rule.SourceList, &l4RuleSource)
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			e := dayuService.ModifyL7Rule(ctx, resourceType, resourceId, rule)
			if e != nil {
				return tccommon.RetryError(e)
			}
			return nil
		})

		if err != nil {
			return err
		}

		readyFlag, rErr := checkL7RuleStatus(meta, resourceType, resourceId, ruleId, "modify")
		if rErr != nil {
			return rErr
		}
		if !readyFlag {
			return fmt.Errorf("Modify operation is timeout %s", ruleId)
		}

	}

	healthFlag := false
	healthKey := []string{"health_check_switch", "health_check_interval", "health_check_path", "health_check_method", "health_check_unhealth_num", "health_check_health_num", "health_check_code"}

	for _, key := range healthKey {
		if d.HasChange(key) {
			healthFlag = true
		}
	}

	if healthFlag {
		//check
		sourceList := d.Get("source_list").(*schema.Set).List()
		if len(sourceList) <= 1 {
			return fmt.Errorf("The `health_check_switch` cannot be set when `source_list` has only one item.")
		}

		//set health check
		var healthCheck dayu.L7HealthConfig
		healthCheck.Protocol = helper.String(d.Get("protocol").(string))
		healthCheck.Domain = &domain
		healthCheck.Enable = helper.BoolToInt64Pointer(d.Get("health_check_switch").(bool))
		healthCheck.Interval = helper.IntUint64(d.Get("health_check_interval").(int))
		healthCheck.Method = helper.String(d.Get("health_check_method").(string))
		healthCheck.Url = helper.String(d.Get("health_check_path").(string))
		healthCheck.KickNum = helper.IntUint64(d.Get("health_check_unhealth_num").(int))
		healthCheck.AliveNum = helper.IntUint64(d.Get("health_check_health_num").(int))
		healthCheck.StatusCode = helper.IntUint64(d.Get("health_check_code").(int))

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			e := dayuService.SetL7Health(ctx, resourceType, resourceId, healthCheck)
			if e != nil {
				return tccommon.RetryError(e)
			}
			return nil
		})

		if err != nil {
			return err
		}

		readyFlag, rErr := checkL7RuleStatus(meta, resourceType, resourceId, ruleId, "check_health")
		if rErr != nil {
			return rErr
		}
		if !readyFlag {
			return fmt.Errorf("Set health is timeout %s", ruleId)
		}
	}

	if d.HasChange("switch") {
		//set switch
		switchFlag := d.Get("switch").(bool)
		protocol := d.Get("protocol").(string)
		if d.HasChange("protocol") {
			//set old protocol para close first
			oldInterface, newInterface := d.GetChange("protocol")
			oldProtocol := oldInterface.(string)
			newProtocol := newInterface.(string)
			protocol = oldProtocol
			//open new only
			if switchFlag {
				protocol = newProtocol
			} else {
				protocol = ""
			}
		}
		if protocol != "" {

			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				e := dayuService.SetRuleSwitch(ctx, resourceType, resourceId, ruleId, switchFlag, protocol)
				if e != nil {
					return tccommon.RetryError(e, tccommon.InternalError)
				}
				return nil
			})

			if err != nil {
				return err
			}

			//check switch status
			readyFlag, rErr := checkL7RuleStatus(meta, resourceType, resourceId, ruleId, fmt.Sprintf("check_switch_%t", switchFlag))
			if rErr != nil {
				return rErr
			}
			if !readyFlag {
				return fmt.Errorf("Set switch is timeout %s", ruleId)
			}
		}

	}

	d.Partial(false)

	return resourceTencentCloudDayuL7RuleRead(d, meta)
}

func resourceTencentCloudDayuL7RuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_l7_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 3 {
		return fmt.Errorf("broken ID of dayu L7 rule")
	}
	resourceType := items[0]
	resourceId := items[1]
	ruleId := items[2]

	dayuService := DayuService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	//set rule
	rule, health, has, err := dayuService.DescribeL7Rule(ctx, resourceType, resourceId, ruleId)
	if err != nil {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			rule, health, has, err = dayuService.DescribeL7Rule(ctx, resourceType, resourceId, ruleId)
			if err != nil {
				return tccommon.RetryError(err)
			}
			return nil
		})
	}
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}
	_ = d.Set("protocol", rule.Protocol)
	_ = d.Set("domain", rule.Domain)
	_ = d.Set("rule_id", rule.RuleId)
	_ = d.Set("ssl_id", rule.SSLId)
	_ = d.Set("name", rule.RuleName)
	_ = d.Set("source_type", int(*rule.SourceType))
	_ = d.Set("status", int(*rule.Status))
	sourceList := make([]*string, 0, len(rule.SourceList))
	for _, v := range rule.SourceList {
		sourceList = append(sourceList, v.Source)
	}
	_ = d.Set("source_list", helper.StringsInterfaces(sourceList))

	if *rule.Protocol == DAYU_L7_RULE_PROTOCOL_HTTPS {
		_ = d.Set("switch", *rule.CCEnable > 0)
	} else {
		_ = d.Set("switch", *rule.CCStatus > 0)
	}
	//set health check
	if health == nil {
		_ = d.Set("health_check_switch", false)
		return nil
	}

	_ = d.Set("health_check_switch", *health.Enable > 0)
	_ = d.Set("health_check_path", health.Url)
	_ = d.Set("health_check_method", health.Method)
	_ = d.Set("health_check_health_num", int(*health.AliveNum))
	_ = d.Set("health_check_unhealth_num", int(*health.KickNum))
	_ = d.Set("health_check_interval", int(*health.Interval))
	_ = d.Set("health_check_code", int(*health.StatusCode))

	return nil
}

func resourceTencentCloudDayuL7RuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_l7_rule.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 3 {
		return fmt.Errorf("broken ID of L7 rule")
	}
	resourceType := items[0]
	resourceId := items[1]
	ruleId := items[2]

	dayuService := DayuService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := dayuService.DeleteL7Rule(ctx, resourceType, resourceId, ruleId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})

	if err != nil {
		return err
	}

	readyFlag, rErr := checkL7RuleStatus(meta, resourceType, resourceId, ruleId, "delete")
	if rErr != nil {
		return rErr
	}
	if !readyFlag {
		return fmt.Errorf("Delete is timeout %s", ruleId)
	}

	return nil
}

func checkL7RuleStatus(meta interface{}, resourceType string, resourceId string, ruleId string, checkType string) (status bool, errRrt error) {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_l7_rule.check_status")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	dayuService := DayuService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		sRule, health, has, err := dayuService.DescribeL7Rule(ctx, resourceType, resourceId, ruleId)
		if err != nil {
			return tccommon.RetryError(err)
		}

		if has {
			//created failed
			if *sRule.Status == DAYU_L7_STATUS_SET_FAIL && (checkType == "create" || checkType == "modify") {
				err = fmt.Errorf("%s rule %s failed...", checkType, ruleId)
				status = false
				return resource.NonRetryableError(err)
			} else if *sRule.Status == DAYU_L7_STATUS_SET_DONE && (checkType == "create" || checkType == "modify") {
				//action completed
				status = true
				return nil
			} else if *sRule.Status == DAYU_L7_STATUS_DEL_FAIL && checkType == "delete" {
				//delete failed
				err = fmt.Errorf("%s rule %s failed...", checkType, ruleId)
				status = false
				return resource.NonRetryableError(err)
			} else if health != nil && *health.Status == DAYU_L7_HEALTH_STATUS_SET_DONE && checkType == "check_health" {
				//check health setting completed
				status = true
				return nil
			} else if health != nil && *health.Status == DAYU_L7_HEALTH_STATUS_SET_FAIL && checkType == "check_health" {
				//check health setting failed
				status = false
				err = fmt.Errorf("%s rule %s failed...status %d", checkType, ruleId, *sRule.Status)
				return resource.NonRetryableError(err)
			} else if checkType == "check_switch_true" {
				//check switch set on completed, the para of http is different from https
				if (*sRule.Protocol == DAYU_L7_RULE_PROTOCOL_HTTPS && *sRule.CCEnable == 1) || (*sRule.Protocol == DAYU_L7_RULE_PROTOCOL_HTTP && *sRule.CCStatus == 1) {
					status = true
					return nil
				} else {
					//check switch set on failed
					status = false
					err = fmt.Errorf("%s rule %s ...", checkType, ruleId)
					return resource.RetryableError(err)
				}
			} else if checkType == "check_switch_false" {
				//check switch set off completed, same as above
				if (*sRule.Protocol == DAYU_L7_RULE_PROTOCOL_HTTPS && *sRule.CCEnable == 0) || (*sRule.Protocol == DAYU_L7_RULE_PROTOCOL_HTTP && *sRule.CCStatus == 0) {
					status = true
					return nil
				} else {
					//check switch set off failed
					status = false
					err = fmt.Errorf("%s rule %s ...", checkType, ruleId)
					return resource.RetryableError(err)
				}
			} else {
				if *sRule.Status == DAYU_L7_STATUS_SSL_WAIT {
					//wait to check ssl
					err = fmt.Errorf("SSL is not ready")
				} else {
					//other cases lead to retry(delete done, set waiting, delete waiting, health setting)
					err = fmt.Errorf("%s rule %s wait to be done, Status %d... describe retry", checkType, ruleId, *sRule.Status)
				}
				return resource.RetryableError(err)
			}
		} else {
			if checkType == "delete" {
				//check delete and do not exist, consider success
				status = true
				return nil
			} else {
				//other cases with no exist, report error
				err = fmt.Errorf("cannot find %s rule", ruleId)
				return resource.NonRetryableError(err)
			}
		}
	})

	if err != nil {
		status = false
	}
	return status, err
}
