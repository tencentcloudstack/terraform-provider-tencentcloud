/*
Use this resource to create a dayu CC self-define http policy

Example Usage

```hcl
resource "tencentcloud_dayu_cc_http_policy" "test_bgpip" {
  resource_type = "bgpip"
  resource_id   = "bgpip-00000294"
  name          = "policy_match"
  smode         = "matching"
  action        = "drop"
  switch        = true
  rule_list {
    skey     = "host"
    operator = "include"
    value    = "123"
  }
}

resource "tencentcloud_dayu_cc_http_policy" "test_net" {
  resource_type = "net"
  resource_id   = "net-0000007e"
  name          = "policy_match"
  smode         = "matching"
  action        = "drop"
  switch        = true
  rule_list {
    skey     = "cgi"
    operator = "equal"
    value    = "123"
  }
}

resource "tencentcloud_dayu_cc_http_policy" "test_bgpmultip" {
  resource_type = "bgp-multip"
  resource_id   = "bgp-0000008o"
  name          = "policy_match"
  smode         = "matching"
  action        = "alg"
  switch        = true
  ip            = "111.230.178.25"

  rule_list {
    skey     = "referer"
    operator = "not_include"
    value    = "123"
  }
}

resource "tencentcloud_dayu_cc_http_policy" "test_bgp" {
  resource_type = "bgp"
  resource_id   = "bgp-000006mq"
  name          = "policy_match"
  smode         = "matching"
  action        = "alg"
  switch        = true

  rule_list {
    skey     = "ua"
    operator = "not_include"
    value    = "123"
  }
}
```
*/
package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	//sdkError "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	dayu "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dayu/v20180709"
)

func resourceTencentCloudDayuCCHttpPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDayuCCHttpPolicyCreate,
		Read:   resourceTencentCloudDayuCCHttpPolicyRead,
		Update: resourceTencentCloudDayuCCHttpPolicyUpdate,
		Delete: resourceTencentCloudDayuCCHttpPolicyDelete,

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the resource that the CC self-define http policy works for.",
			},
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue(DAYU_RESOURCE_TYPE),
				ForceNew:     true,
				Description:  "Type of the resource that the CC self-define http policy works for, valid values are `bgpip`, `bgp`, `bgp-multip` and `net`.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateStringLengthInRange(1, 20),
				Description:  "Name of the CC self-define http policy. Length should between 1 and 20.",
			},
			"smode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(DAYU_CC_POLICY_SMODE),
				Default:      DAYU_CC_POLICY_SMODE_MATCH,
				Description:  "Match mode, and valid values are `matching`, `speedlimit`. Note: the speed limit type CC self-define policy can only set one.",
			},
			"frequency": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIntegerInRange(1, 10000),
				Description:  "Max frequency per minute, only valid when `smode` is `speedlimit`, the valid value ranges from 1 to 10000.",
			},
			"action": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateAllowedStringValue(DAYU_CC_POLICY_ACTION),
				Description:  "Action mode, only valid when `smode` is `matching`. Valid values are `alg` and `drop`.",
			},
			"switch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicate the CC self-define http policy takes effect or not.",
			},
			"rule_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"skey": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "host",
							ValidateFunc: validateAllowedStringValue(DAYU_CC_POLICY_HTTP_CHECK_TYPE),
							Description:  "Key of the rule, valid values are `host`, `cgi`, `ua`, `referer`.",
						},
						"operator": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "include",
							ValidateFunc: validateAllowedStringValue(DAYU_CC_POLICY_CHECK_OP),
							Description:  "Operator of the rule, valid values are `include`, `not_include`, `equal`.",
						},
						"value": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "",
							ValidateFunc: validateStringLengthInRange(0, 31),
							Description:  "Rule value, then length should be less than 31 bytes.",
						},
					},
				},
				Description: "Rule list of the CC self-define http policy,  only valid when `smode` is `matching`.",
			},
			"ip": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validateIp,
				Description:  "Ip of the CC self-define http policy, only valid when `resource_type` is `bgp-multip`. The num of list items can only be set one.",
			},
			//computed
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the CC self-define http policy.",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Id of the CC self-define http policy.",
			},
		},
	}
}

func resourceTencentCloudDayuCCHttpPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dayu_cc_http_policy.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	resourceId := d.Get("resource_id").(string)
	resourceType := d.Get("resource_type").(string)
	//set CCPolicy
	var ccPolicy dayu.CCPolicy
	ccPolicy.Name = helper.String(d.Get("name").(string))
	smode := d.Get("smode").(string)
	ccPolicy.Smode = &smode
	frequency := 0
	if v, ok := d.GetOk("frequency"); ok {
		frequency = v.(int)
	}

	if smode == DAYU_CC_POLICY_SMODE_SPEED_LIMIT {
		if frequency == 0 {
			return fmt.Errorf("`frequencys` should be set when `smode` is `speedlimit`.")
		}
		ccPolicy.Frequency = helper.IntUint64(frequency)
	} else {
		ccPolicy.ExeMode = helper.String(d.Get("action").(string))
	}
	ccPolicy.Protocol = helper.String(DAYU_L7_RULE_PROTOCOL_HTTP)
	switchFlag := d.Get("switch").(bool)
	if switchFlag {
		ccPolicy.Switch = helper.IntUint64(1)
	} else {
		ccPolicy.Switch = helper.IntUint64(0)
	}

	ip := ""
	if v, ok := d.GetOk("ip"); ok {
		ip = v.(string)
	}
	if ip != "" {
		ccPolicy.IpList = []*string{&ip}
	} else {
		ccPolicy.IpList = []*string{}
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

	dayuService := DayuService{client: meta.(*TencentCloudClient).apiV3Conn}

	policyId := ""
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := dayuService.CreateCCSelfdefinePolicy(ctx, resourceType, resourceId, ccPolicy)
		if e != nil {
			return retryError(e)
		}
		policyId = result
		return nil
	})

	if err != nil {
		return err
	}

	d.SetId(resourceType + FILED_SP + resourceId + FILED_SP + policyId)

	return resourceTencentCloudDayuCCHttpPolicyRead(d, meta)
}

func resourceTencentCloudDayuCCHttpPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dayu_cc_http_policy.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	items := strings.Split(d.Id(), FILED_SP)
	if len(items) < 3 {
		return fmt.Errorf("broken ID of dayu CC policy")
	}
	resourceType := items[0]
	resourceId := items[1]
	policyId := items[2]

	dayuService := DayuService{client: meta.(*TencentCloudClient).apiV3Conn}

	policy, has, err := dayuService.DescribeCCSelfdefinePolicy(ctx, resourceType, resourceId, policyId)
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			policy, has, err = dayuService.DescribeCCSelfdefinePolicy(ctx, resourceType, resourceId, policyId)
			if err != nil {
				return retryError(err)
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
	_ = d.Set("smode", policy.Smode)
	_ = d.Set("policy_id", policy.SetId)
	_ = d.Set("action", policy.ExeMode)
	ipList := helper.StringsInterfaces(policy.IpList)
	if len(ipList) == 1 {
		_ = d.Set("ip", ipList[0])
	}
	_ = d.Set("switch", *policy.Switch > 0)

	if policy.Frequency != nil && *policy.Smode == "frequency" {
		_ = d.Set("frequency", policy.Frequency)
	}
	if policy.RuleList != nil && *policy.Smode == "matching" {
		_ = d.Set("rule_list", flattenCCRuleList(policy.RuleList))
	}

	return nil
}

func resourceTencentCloudDayuCCHttpPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dayu_cc_http_policy.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	items := strings.Split(d.Id(), FILED_SP)
	if len(items) < 3 {
		return fmt.Errorf("broken ID of CC policy")
	}
	resourceType := items[0]
	resourceId := items[1]
	policyId := items[2]

	//set CCPolicy
	var ccPolicy dayu.CCPolicy
	ccPolicy.Name = helper.String(d.Get("name").(string))
	smode := d.Get("smode").(string)
	ccPolicy.Smode = &smode
	frequency := 0
	if v, ok := d.GetOk("frequency"); ok {
		frequency = v.(int)
	}
	if smode == DAYU_CC_POLICY_SMODE_SPEED_LIMIT {
		if frequency == 0 {
			return fmt.Errorf("`speedlimit` should be set when `smode` is `speedlimit`.")
		}
		ccPolicy.Frequency = helper.IntUint64(frequency)
	} else {
		ccPolicy.ExeMode = helper.String(d.Get("action").(string))
	}
	ccPolicy.Protocol = helper.String(DAYU_L7_RULE_PROTOCOL_HTTP)

	switchFlag := d.Get("switch").(bool)
	if switchFlag {
		ccPolicy.Switch = helper.IntUint64(1)
	} else {
		ccPolicy.Switch = helper.IntUint64(0)
	}

	ruleList := d.Get("rule_list").([]interface{})
	ccPolicy.RuleList = make([]*dayu.CCRule, 0, len(ruleList))
	for _, rule := range ruleList {
		var ccRule dayu.CCRule
		ruleMap := rule.(map[string]interface{})
		ccRule.Skey = helper.String(ruleMap["skey"].(string))
		ccRule.Operator = helper.String(ruleMap["operator"].(string))
		ccRule.Value = helper.String(ruleMap["value"].(string))
		ccPolicy.RuleList = append(ccPolicy.RuleList, &ccRule)
	}
	ip := ""
	if v, ok := d.GetOk("ip"); ok {
		ip = v.(string)
	}
	if ip != "" {
		ccPolicy.IpList = []*string{&ip}
	} else {
		ccPolicy.IpList = []*string{}
	}
	dayuService := DayuService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := dayuService.ModifyCCSelfdefinePolicy(ctx, resourceType, resourceId, policyId, ccPolicy)
		if e != nil {
			return retryError(e)
		}
		return nil
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudDayuCCHttpPolicyRead(d, meta)
}

func resourceTencentCloudDayuCCHttpPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dayu_cc_http_policy.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	items := strings.Split(d.Id(), FILED_SP)
	if len(items) < 3 {
		return fmt.Errorf("broken ID of CC policy")
	}
	resourceType := items[0]
	resourceId := items[1]
	policyId := items[2]

	dayuService := DayuService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := dayuService.DeleteCCSelfdefinePolicy(ctx, resourceType, resourceId, policyId)
		if e != nil {
			return retryError(e)
		}
		return nil
	})

	if err != nil {
		return err
	}

	_, has, err := dayuService.DescribeCCSelfdefinePolicy(ctx, resourceType, resourceId, policyId)
	if err != nil || has {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, has, err = dayuService.DescribeCCSelfdefinePolicy(ctx, resourceType, resourceId, policyId)
			if err != nil {
				return retryError(err)
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
