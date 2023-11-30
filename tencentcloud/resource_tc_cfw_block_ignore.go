package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfw "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCfwBlockIgnore() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCfwBlockIgnoreCreate,
		Read:   resourceTencentCloudCfwBlockIgnoreRead,
		Update: resourceTencentCloudCfwBlockIgnoreUpdate,
		Delete: resourceTencentCloudCfwBlockIgnoreDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"domain"},
				Description:  "Rule IP address, one of IP and Domain is required.",
			},
			"domain": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"ip"},
				Description:  "Rule domain name, one of IP and Domain is required.",
			},
			"direction": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue(DIRECTION),
				Description:  "Rule direction, 0 outbound, 1 inbound, 3 intranet.",
			},
			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule end time, format: 2006-01-02 15:04:05, must be greater than the current time.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Remarks information, length cannot exceed 50.",
			},
			"start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Rule start time.",
			},
			"rule_type": {
				Required:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validateAllowedIntValue(RULE_TYPE),
				Description:  "Rule type, 1 block, 2 ignore, domain block is not supported.",
			},
		},
	}
}

func resourceTencentCloudCfwBlockIgnoreCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_block_ignore.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId                = getLogId(contextNil)
		request              = cfw.NewCreateBlockIgnoreRuleListRequest()
		intrusionDefenseRule = cfw.IntrusionDefenseRule{}
		iP                   string
		domain               string
		direction            string
		ruleType             string
	)

	if v, ok := d.GetOk("ip"); ok {
		intrusionDefenseRule.IP = helper.String(v.(string))
		iP = v.(string)
	}

	if v, ok := d.GetOk("domain"); ok {
		intrusionDefenseRule.Domain = helper.String(v.(string))
		domain = v.(string)
	}

	if v, ok := d.GetOk("direction"); ok {
		directionInt, _ := strconv.ParseInt(v.(string), 10, 64)
		intrusionDefenseRule.Direction = &directionInt
		direction = v.(string)
	}

	if v, ok := d.GetOk("end_time"); ok {
		intrusionDefenseRule.EndTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("comment"); ok {
		intrusionDefenseRule.Comment = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_time"); ok {
		intrusionDefenseRule.StartTime = helper.String(v.(string))
	}

	request.Rules = append(request.Rules, &intrusionDefenseRule)

	if v, ok := d.GetOkExists("rule_type"); ok {
		request.RuleType = helper.IntInt64(v.(int))
		ruleType = strconv.Itoa(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCfwClient().CreateBlockIgnoreRuleList(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cfw blockIgnore failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{iP, domain, direction, ruleType}, FILED_SP))

	return resourceTencentCloudCfwBlockIgnoreRead(d, meta)
}

func resourceTencentCloudCfwBlockIgnoreRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_block_ignore.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = CfwService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	iP := idSplit[0]
	domain := idSplit[1]
	direction := idSplit[2]
	ruleType := idSplit[3]

	blockIgnoreRule, err := service.DescribeCfwBlockIgnoreListById(ctx, iP, domain, direction, ruleType)
	if err != nil {
		return err
	}

	if blockIgnoreRule == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CfwBlockIgnore` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if blockIgnoreRule.Ioc != nil {
		_ = d.Set("ip", blockIgnoreRule.Ioc)
	}

	if blockIgnoreRule.Domain != nil {
		_ = d.Set("domain", blockIgnoreRule.Domain)
	}

	if blockIgnoreRule.Direction != nil {
		directionStr := strconv.FormatInt(*blockIgnoreRule.Direction, 10)
		_ = d.Set("direction", directionStr)
	}

	if blockIgnoreRule.EndTime != nil {
		_ = d.Set("end_time", blockIgnoreRule.EndTime)
	}

	if blockIgnoreRule.Comment != nil {
		_ = d.Set("comment", blockIgnoreRule.Comment)
	}

	if blockIgnoreRule.StartTime != nil {
		_ = d.Set("start_time", blockIgnoreRule.StartTime)
	}

	if blockIgnoreRule.Action != nil {
		_ = d.Set("rule_type", blockIgnoreRule.Action)
	}

	return nil
}

func resourceTencentCloudCfwBlockIgnoreUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_block_ignore.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId                = getLogId(contextNil)
		request              = cfw.NewModifyBlockIgnoreRuleRequest()
		intrusionDefenseRule = cfw.IntrusionDefenseRule{}
	)

	immutableArgs := []string{"ip", "domain", "direction", "rule_type"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	iP := idSplit[0]
	domain := idSplit[1]
	direction := idSplit[2]
	ruleType := idSplit[3]

	if iP != "" {
		intrusionDefenseRule.IP = &iP
	} else {
		intrusionDefenseRule.Domain = &domain
	}

	directionInt, _ := strconv.ParseInt(direction, 10, 64)
	ruleTypeInt, _ := strconv.ParseInt(ruleType, 10, 64)
	intrusionDefenseRule.Direction = &directionInt
	request.RuleType = &ruleTypeInt

	if v, ok := d.GetOk("end_time"); ok {
		intrusionDefenseRule.EndTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("comment"); ok {
		intrusionDefenseRule.Comment = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_time"); ok {
		intrusionDefenseRule.StartTime = helper.String(v.(string))
	}

	request.Rule = &intrusionDefenseRule

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCfwClient().ModifyBlockIgnoreRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update cfw blockIgnore failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCfwBlockIgnoreRead(d, meta)
}

func resourceTencentCloudCfwBlockIgnoreDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_block_ignore.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = CfwService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	iP := idSplit[0]
	domain := idSplit[1]
	direction := idSplit[2]
	ruleType := idSplit[3]

	if err := service.DeleteCfwBlockIgnoreListById(ctx, iP, domain, direction, ruleType); err != nil {
		return err
	}

	return nil
}
