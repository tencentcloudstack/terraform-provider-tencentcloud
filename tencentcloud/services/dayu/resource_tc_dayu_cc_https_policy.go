package dayu

import (
	"context"
	"errors"
	"fmt"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dayu "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dayu/v20180709"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDayuCCHttpsPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDayuCCHttpsPolicyCreate,
		Read:   resourceTencentCloudDayuCCHttpsPolicyRead,
		Update: resourceTencentCloudDayuCCHttpsPolicyUpdate,
		Delete: resourceTencentCloudDayuCCHttpsPolicyDelete,

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the resource that the CC self-define https policy works for.",
			},
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_RESOURCE_TYPE_HTTPS),
				ForceNew:     true,
				Description:  "Type of the resource that the CC self-define https policy works for, valid value is `bgpip`.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 20),
				Description:  "Name of the CC self-define https policy. Length should between 1 and 20.",
			},
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Domain that the CC self-define https policy works for, only valid when `protocol` is `https`.",
			},
			"rule_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Rule id of the domain that the CC self-define https policy works for, only valid when `protocol` is `https`.",
			},
			"switch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicate the CC self-define https policy takes effect or not.",
			},
			"action": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_CC_POLICY_ACTION),
				Description:  "Action mode. Valid values are `alg` and `drop`.",
			},
			"rule_list": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"skey": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_CC_POLICY_HTTPS_CHECK_TYPE),
							Description:  "Key of the rule. Valid values are `cgi`, `ua` and `referer`.",
						},
						"operator": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_CC_POLICY_CHECK_OP_HTTPS),
							Description:  "Operator of the rule. Valid values are `include` and `equal`.",
						},
						"value": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: tccommon.ValidateStringLengthInRange(0, 31),
							Description:  "Rule value, then length should be less than 31 bytes.",
						},
					},
				},
				Description: "Rule list of the CC self-define https policy.",
			},
			//computed
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the CC self-define https policy.",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Id of the CC self-define https policy.",
			},
			"ip_list": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Ip of the CC self-define https policy.",
			},
		},
	}
}

func resourceTencentCloudDayuCCHttpsPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_cc_https_policy.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	resourceId := d.Get("resource_id").(string)
	resourceType := d.Get("resource_type").(string)
	//set CCPolicy
	var ccPolicy dayu.CCPolicy
	ccPolicy.Name = helper.String(d.Get("name").(string))
	ccPolicy.Smode = helper.String(DAYU_CC_POLICY_SMODE_MATCH)
	ccPolicy.Protocol = helper.String(DAYU_L7_RULE_PROTOCOL_HTTPS)
	ccPolicy.Domain = helper.String(d.Get("domain").(string))
	ccPolicy.RuleId = helper.String(d.Get("rule_id").(string))
	ccPolicy.ExeMode = helper.String(d.Get("action").(string))

	ccPolicy.IpList = []*string{}

	if v, ok := d.GetOk("rule_id"); ok {
		ccPolicy.RuleId = helper.String(v.(string))
	}

	switchFlag := d.Get("switch").(bool)
	if switchFlag {
		ccPolicy.Switch = helper.IntUint64(1)
	} else {
		ccPolicy.Switch = helper.IntUint64(0)
	}

	ruleList := d.Get("rule_list").(*schema.Set).List()
	ccPolicy.RuleList = make([]*dayu.CCRule, 0, len(ruleList))
	for _, rule := range ruleList {
		var ccRule dayu.CCRule
		ruleMap := rule.(map[string]interface{})
		ccRule.Skey = helper.String(ruleMap["skey"].(string))
		ccRule.Operator = helper.String(ruleMap["operator"].(string))
		ccRule.Value = helper.String(ruleMap["value"].(string))
		ccPolicy.RuleList = append(ccPolicy.RuleList, &ccRule)
	}

	dayuService := DayuService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	policyId := ""

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := dayuService.CreateCCSelfdefinePolicy(ctx, resourceType, resourceId, ccPolicy)
		if e != nil {
			return tccommon.RetryError(e)
		}
		policyId = result
		return nil
	})

	if err != nil {
		return err
	}

	d.SetId(resourceType + tccommon.FILED_SP + resourceId + tccommon.FILED_SP + policyId)

	return resourceTencentCloudDayuCCHttpPolicyRead(d, meta)
}

func resourceTencentCloudDayuCCHttpsPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_cc_https_policy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 3 {
		return fmt.Errorf("broken ID of dayu CC policy")
	}
	resourceType := items[0]
	resourceId := items[1]
	policyId := items[2]

	dayuService := DayuService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	policy, has, err := dayuService.DescribeCCSelfdefinePolicy(ctx, resourceType, resourceId, policyId)
	if err != nil {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			policy, has, err = dayuService.DescribeCCSelfdefinePolicy(ctx, resourceType, resourceId, policyId)
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
	_ = d.Set("name", policy.Name)
	_ = d.Set("create_time", policy.CreateTime)
	_ = d.Set("policy_id", policy.SetId)
	_ = d.Set("action", policy.ExeMode)
	_ = d.Set("rule_id", policy.RuleId)
	_ = d.Set("domain", policy.Domain)
	_ = d.Set("switch", *policy.Switch > 0)
	_ = d.Set("rule_list", flattenCCRuleList(policy.RuleList))
	_ = d.Set("ip_list", helper.StringsInterfaces(policy.IpList))

	return nil
}

func resourceTencentCloudDayuCCHttpsPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_cc_https_policy.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 3 {
		return fmt.Errorf("broken ID of CC policy")
	}
	resourceType := items[0]
	resourceId := items[1]
	policyId := items[2]

	//set CCPolicy
	var ccPolicy dayu.CCPolicy
	ccPolicy.Name = helper.String(d.Get("name").(string))
	ccPolicy.Smode = helper.String(DAYU_CC_POLICY_SMODE_MATCH)
	ccPolicy.Protocol = helper.String(DAYU_L7_RULE_PROTOCOL_HTTP)
	ccPolicy.SetId = helper.String(policyId)

	if v, ok := d.GetOk("rule_id"); ok {
		ccPolicy.RuleId = helper.String(v.(string))
	}

	switchFlag := d.Get("switch").(bool)
	if switchFlag {
		ccPolicy.Switch = helper.IntUint64(1)
	} else {
		ccPolicy.Switch = helper.IntUint64(0)
	}

	ccPolicy.ExeMode = helper.String(d.Get("action").(string))
	ruleList := d.Get("rule_list").(*schema.Set).List()
	ccPolicy.RuleList = make([]*dayu.CCRule, 0, len(ruleList))
	for _, rule := range ruleList {
		var ccRule dayu.CCRule
		ruleMap := rule.(map[string]interface{})
		ccRule.Skey = helper.String(ruleMap["skey"].(string))
		ccRule.Operator = helper.String(ruleMap["operator"].(string))
		ccRule.Value = helper.String(ruleMap["value"].(string))
		ccPolicy.RuleList = append(ccPolicy.RuleList, &ccRule)
	}

	//the sdk really designed error, it need this para
	ipList := d.Get("ip_list").(*schema.Set).List()
	for _, ip := range ipList {
		ccPolicy.IpList = append(ccPolicy.IpList, helper.String(ip.(string)))
	}
	dayuService := DayuService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := dayuService.ModifyCCSelfdefinePolicy(ctx, resourceType, resourceId, policyId, ccPolicy)
		if e != nil {
			return tccommon.RetryError(e, tccommon.InternalError)
		}
		return nil
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudDayuCCHttpPolicyRead(d, meta)
}

func resourceTencentCloudDayuCCHttpsPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_cc_https_policy.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 3 {
		return fmt.Errorf("broken ID of CC policy")
	}
	resourceType := items[0]
	resourceId := items[1]
	policyId := items[2]

	dayuService := DayuService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := dayuService.DeleteCCSelfdefinePolicy(ctx, resourceType, resourceId, policyId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})

	if err != nil {
		return err
	}

	_, has, err := dayuService.DescribeCCSelfdefinePolicy(ctx, resourceType, resourceId, policyId)
	if err != nil || has {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, has, err = dayuService.DescribeCCSelfdefinePolicy(ctx, resourceType, resourceId, policyId)
			if err != nil {
				return tccommon.RetryError(err)
			}

			if has {
				err = fmt.Errorf("delete DDoS policy fail, CC policy still exist from sdk DescribeCCSelfDefinePolicy")
				return resource.RetryableError(err)
			}

			return nil
		})
	}
	if err != nil {
		return err
	}
	if !has {
		return nil
	} else {
		return errors.New("delete CC policy fail, CC policy still exist from sdk DescribeCCSelfDefinePolicy")
	}
}
