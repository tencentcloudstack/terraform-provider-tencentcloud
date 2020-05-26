/*
Use this resource to create dayu layer 4 rule

~> **NOTE:** This resource only support resource Anti-DDoS of type `bgpip` and `net`

Example Usage

```hcl
resource "tencentcloud_dayu_l4_rule" "test_rule" {
  resource_type             = "bgpip"
  resource_id               = "bgpip-00000294"
  name                      = "rule_test"
  protocol                  = "TCP"
  s_port                    = 80
  d_port                    = 60
  source_type               = 2
  health_check_switch       = true
  health_check_timeout      = 30
  health_check_interval     = 35
  health_check_health_num   = 5
  health_check_unhealth_num = 10
  session_switch            = false
  session_time              = 300

  source_list {
    source = "1.1.1.1"
    weight = 100
  }
  source_list {
    source = "2.2.2.2"
    weight = 50
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
	dayu "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dayu/v20180709"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDayuL4Rule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDayuL4RuleCreate,
		Read:   resourceTencentCloudDayuL4RuleRead,
		Update: resourceTencentCloudDayuL4RuleUpdate,
		Delete: resourceTencentCloudDayuL4RuleDelete,

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the resource that the layer 4 rule works for.",
			},
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue(DAYU_RESOURCE_TYPE_RULE),
				ForceNew:     true,
				Description:  "Type of the resource that the layer 4 rule works for, valid values are `bgpip` and `net`.",
			},
			"source_type": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedIntValue(DAYU_L7_RULE_SOURCE_TYPE),
				Description:  "Source type, 1 for source of host, 2 for source of ip.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the rule. When the `resource_type` is `net`, this field should be set with valid domain.",
			},
			"s_port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validatePort,
				Description:  "The source port of the L4 rule.",
			},
			"d_port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validatePort,
				Description:  "The destination port of the L4 rule.",
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue(DAYU_L4_RULE_PROTOCOL),
				Description:  "Protocol of the rule, valid values are `http`, `https`. When `source_type` is 1(host source), the value of this field can only set with `tcp`.",
			},
			"source_list": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Source ip or domain, valid format of ip is like `1.1.1.1` and valid format of host source is like `abc.com`.",
						},
						"weight": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validateIntegerInRange(0, 100),
							Description:  "Weight of the source, the valid value ranges from 0 to 100.",
						},
					},
				},
				MinItems:    1,
				MaxItems:    20,
				Description: "Source list of the rule, it can be a set of ip sources or a set of domain sources. The number of items ranges from 1 to 20.",
			},
			"health_check_switch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates whether health check is enabled. The default is `false`. Only valid when source list has more than one source item.",
			},
			"health_check_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIntegerInRange(10, 60),
				Description:  "Interval time of health check. The value range is 10-60 sec, and the default is 15 sec.",
			},
			"health_check_health_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIntegerInRange(2, 10),
				Description:  "Health threshold of health check, and the default is 3. If a success result is returned for the health check 3 consecutive times, indicates that the forwarding is normal. The value range is 2-10.",
			},
			"health_check_unhealth_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIntegerInRange(2, 10),
				Description:  "Unhealthy threshold of health check, and the default is 3. If the unhealthy result is returned 3 consecutive times, indicates that the forwarding is abnormal. The value range is 2-10.",
			},
			"health_check_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIntegerInRange(2, 60),
				Description:  "HTTP Status Code. The default is 26 and value range is 2-60.",
			},
			"session_switch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicate that the session will keep or not, and default value is `false`.",
			},
			"session_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIntegerInRange(1, 300),
				Description:  "Session keep time, only valid when `session_switch` is true, the available value ranges from 1 to 300 and unit is second.",
			},
			//computed
			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Id of the layer 4 rule.",
			},
			"lb_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "LB type of the rule, 1 for weight cycling and 2 for IP hash.",
			},
		},
	}
}

func resourceTencentCloudDayuL4RuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dayu_l4_rule.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	resourceId := d.Get("resource_id").(string)
	resourceType := d.Get("resource_type").(string)

	destPort := d.Get("d_port").(int)

	//check
	protocol := d.Get("protocol").(string)
	source_type := d.Get("source_type").(int)
	sourceList := d.Get("source_list").(*schema.Set).List()

	if source_type == DAYU_L7_RULE_SOURCE_TYPE_HOST && protocol != DAYU_L4_RULE_PROTOCOL_TCP {
		return fmt.Errorf("`protocol` can only be set with `TCP` when `source_type` is 1(host source).")
	}

	healthCheckSwitch := d.Get("health_check_switch").(bool)
	if healthCheckSwitch {
		if len(sourceList) <= 1 {
			return fmt.Errorf("The `health_check_switch` cannot be set `true` when `source_list` has only one item.")
		}
	}

	//check
	timeout := 0
	interval := 0
	if v, ok := d.GetOk("health_check_timeout"); ok {
		timeout = v.(int)
	}
	if v, ok := d.GetOk("health_check_interval"); ok {
		interval = v.(int)
	}

	if timeout > 0 && interval > 0 {
		if timeout > interval {
			return fmt.Errorf("The `health_check_interval` should be greater than `health_check_timeout`.")
		}
	}

	//set L4RuleEntry
	var rule dayu.L4RuleEntry
	rule.LbType = helper.IntUint64(1)
	rule.SourcePort = helper.IntUint64(d.Get("s_port").(int))
	rule.VirtualPort = helper.IntUint64(destPort)
	rule.Protocol = helper.String(d.Get("protocol").(string))
	rule.RuleName = helper.String(d.Get("name").(string))
	sourceType := d.Get("source_type").(int)
	//check that there is no check with the source list and sdk returns
	rule.SourceType = helper.IntUint64(sourceType)

	rule.SourceList = make([]*dayu.L4RuleSource, 0, len(sourceList))
	for _, source := range sourceList {
		sourceMap := source.(map[string]interface{})
		var l4RuleSource dayu.L4RuleSource
		l4RuleSource.Source = helper.String(sourceMap["source"].(string))
		l4RuleSource.Weight = helper.IntUint64(sourceMap["weight"].(int))
		rule.SourceList = append(rule.SourceList, &l4RuleSource)
	}

	dayuService := DayuService{client: meta.(*TencentCloudClient).apiV3Conn}

	ruleId := ""
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := dayuService.CreateL4Rule(ctx, resourceType, resourceId, rule)
		if e != nil {
			return retryError(e)
		}
		ruleId = result
		return nil
	})

	if err != nil {
		return err
	}

	d.SetId(resourceType + FILED_SP + resourceId + FILED_SP + ruleId)

	//set health check
	var healthCheck dayu.L4HealthConfig
	healthCheck.Protocol = helper.String(d.Get("protocol").(string))
	healthCheck.Enable = helper.BoolToInt64Pointer(d.Get("health_check_switch").(bool))
	healthCheck.Interval = helper.IntUint64(d.Get("health_check_interval").(int))
	healthCheck.KickNum = helper.IntUint64(d.Get("health_check_unhealth_num").(int))
	healthCheck.AliveNum = helper.IntUint64(d.Get("health_check_health_num").(int))
	healthCheck.TimeOut = helper.IntUint64(d.Get("health_check_timeout").(int))
	healthCheck.VirtualPort = helper.IntUint64(destPort)

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := dayuService.SetL4Health(ctx, resourceType, resourceId, healthCheck)
		if e != nil {
			return retryError(e)
		}
		return nil
	})

	if err != nil {
		return err
	}

	//set session
	sessionFlag := d.Get("session_switch").(bool)
	sessionTime := 0
	if v, ok := d.GetOk("session_time"); ok {
		sessionTime = v.(int)
	}
	if sessionTime == 0 && sessionFlag {
		return fmt.Errorf("`session_time` should be set when `session_switch` is true.")
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := dayuService.SetSession(ctx, resourceType, resourceId, ruleId, sessionFlag, sessionTime)
		if e != nil {
			return retryError(e)
		}
		return nil
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudDayuL4RuleRead(d, meta)
}

func resourceTencentCloudDayuL4RuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dayu_l4_rule.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	items := strings.Split(d.Id(), FILED_SP)
	if len(items) < 3 {
		return fmt.Errorf("broken ID of dayu L4 rule")
	}
	resourceType := items[0]
	resourceId := items[1]
	ruleId := items[2]

	d.Partial(true)

	sourceType := d.Get("source_type").(int)
	protocol := d.Get("protocol").(string)
	//check
	if sourceType == 1 && protocol != DAYU_L4_RULE_PROTOCOL_TCP {
		return fmt.Errorf("`protocol` can only be set with `TCP` when `source_type` is 1(host source).")
	}
	sourceList := d.Get("source_list").(*schema.Set).List()

	//check
	timeout := 0
	interval := 0
	if v, ok := d.GetOk("health_check_timeout"); ok {
		timeout = v.(int)
	}
	if v, ok := d.GetOk("health_check_interval"); ok {
		interval = v.(int)
	}

	if timeout > 0 && interval > 0 {
		if timeout > interval {
			return fmt.Errorf("The `health_check_interval` should be greater than `health_check_timeout`.")
		}
	}

	dayuService := DayuService{client: meta.(*TencentCloudClient).apiV3Conn}

	ruleFlag := false
	ruleKey := []string{"s_port", "d_port", "protocol", "source_list"}

	for _, key := range ruleKey {
		if d.HasChange(key) {
			ruleFlag = true
		}
	}

	if ruleFlag {
		//set L4RuleEntry
		var rule dayu.L4RuleEntry
		rule.LbType = helper.IntUint64(1)
		rule.SourcePort = helper.IntUint64(d.Get("s_port").(int))
		rule.VirtualPort = helper.IntUint64(d.Get("d_port").(int))
		rule.Protocol = helper.String(d.Get("protocol").(string))
		rule.RuleName = helper.String(d.Get("name").(string))
		rule.RuleId = &ruleId

		rule.SourceType = helper.IntUint64(sourceType)
		rule.Protocol = &protocol

		rule.SourceList = make([]*dayu.L4RuleSource, 0, len(sourceList))
		for _, source := range sourceList {
			sourceMap := source.(map[string]interface{})
			var l4RuleSource dayu.L4RuleSource
			l4RuleSource.Source = helper.String(sourceMap["source"].(string))
			l4RuleSource.Weight = helper.IntUint64(sourceMap["weight"].(int))
			rule.SourceList = append(rule.SourceList, &l4RuleSource)
		}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := dayuService.ModifyL4Rule(ctx, resourceType, resourceId, rule)
			if e != nil {
				return retryError(e)
			}
			return nil
		})

		if err != nil {
			return err
		}

		for _, key := range ruleKey {
			d.SetPartial(key)
		}
	}

	healthFlag := false
	healthKey := []string{"health_check_switch", "health_check_interval", "health_check_timeout", "health_check_unhealth_num", "health_check_health_num", "d_port"}

	for _, key := range healthKey {
		if d.HasChange(key) {
			healthFlag = true
		}
	}

	if healthFlag {
		//set health check
		if len(sourceList) <= 1 {
			return fmt.Errorf("The `health_check_switch` cannot be set when `source_list` has only one item.")
		}

		var healthCheck dayu.L4HealthConfig
		healthCheck.Protocol = helper.String(d.Get("protocol").(string))
		healthCheck.Enable = helper.BoolToInt64Pointer(d.Get("health_check_switch").(bool))
		healthCheck.Interval = helper.IntUint64(d.Get("health_check_interval").(int))
		healthCheck.KickNum = helper.IntUint64(d.Get("health_check_unhealth_num").(int))
		healthCheck.AliveNum = helper.IntUint64(d.Get("health_check_health_num").(int))
		healthCheck.TimeOut = helper.IntUint64(d.Get("health_check_timeout").(int))
		healthCheck.VirtualPort = helper.IntUint64(d.Get("d_port").(int))

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := dayuService.SetL4Health(ctx, resourceType, resourceId, healthCheck)
			if e != nil {
				return retryError(e)
			}
			return nil
		})

		if err != nil {
			return err
		}

		for _, key := range healthKey {
			d.SetPartial(key)
		}
	}

	if d.HasChange("session_switch") || d.HasChange("session_time") {
		sessionFlag := d.Get("session_switch").(bool)
		sessionTime := 0
		if v, ok := d.GetOk("session_time"); ok {
			sessionTime = v.(int)
		}
		if sessionTime == 0 && sessionFlag {
			return fmt.Errorf("`session_time` should be set when `session_switch` is true.")
		}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := dayuService.SetSession(ctx, resourceType, resourceId, ruleId, sessionFlag, sessionTime)
			if e != nil {
				return retryError(e)
			}
			return nil
		})

		if err != nil {
			return err
		}
		d.SetPartial("session_switch")
		d.SetPartial("session_time")
	}

	d.Partial(false)

	return resourceTencentCloudDayuL4RuleRead(d, meta)
}

func resourceTencentCloudDayuL4RuleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dayu_l4_rule.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	items := strings.Split(d.Id(), FILED_SP)
	if len(items) < 3 {
		return fmt.Errorf("broken ID of dayu L4 rule")
	}
	resourceType := items[0]
	resourceId := items[1]
	ruleId := items[2]

	dayuService := DayuService{client: meta.(*TencentCloudClient).apiV3Conn}

	//set rule
	rule, health, has, err := dayuService.DescribeL4Rule(ctx, resourceType, resourceId, ruleId)
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			rule, health, has, err = dayuService.DescribeL4Rule(ctx, resourceType, resourceId, ruleId)
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
	_ = d.Set("protocol", rule.Protocol)
	_ = d.Set("s_port", int(*rule.SourcePort))
	_ = d.Set("d_port", int(*rule.VirtualPort))
	_ = d.Set("rule_id", rule.RuleId)
	_ = d.Set("lb_type", int(*rule.LbType))
	_ = d.Set("name", rule.RuleName)
	_ = d.Set("source_type", int(*rule.SourceType))
	_ = d.Set("session_time", int(*rule.KeepTime))
	_ = d.Set("session_switch", *rule.KeepEnable > 0)
	_ = d.Set("source_list", flattenSourceList(rule.SourceList))

	//set health check
	if health == nil {
		_ = d.Set("health_check_switch", false)
		return nil
	}
	_ = d.Set("health_check_switch", *health.Enable > 0)
	_ = d.Set("health_check_timeout", int(*health.TimeOut))
	_ = d.Set("health_check_health_num", int(*health.AliveNum))
	_ = d.Set("health_check_unhealth_num", int(*health.KickNum))
	_ = d.Set("health_check_interval", int(*health.Interval))

	return nil
}

func resourceTencentCloudDayuL4RuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dayu_l4_rule.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	items := strings.Split(d.Id(), FILED_SP)
	if len(items) < 3 {
		return fmt.Errorf("broken ID of L4 rule")
	}
	resourceType := items[0]
	resourceId := items[1]
	ruleId := items[2]

	dayuService := DayuService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := dayuService.DeleteL4Rule(ctx, resourceType, resourceId, ruleId)
		if e != nil {
			return retryError(e)
		}
		return nil
	})

	if err != nil {
		return err
	}

	_, _, has, err := dayuService.DescribeL4Rule(ctx, resourceType, resourceId, ruleId)
	if err != nil || has {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, _, has, err = dayuService.DescribeL4Rule(ctx, resourceType, resourceId, ruleId)
			if err != nil {
				return retryError(err)
			}

			if has {
				err = fmt.Errorf("delete L4 rule fail, L4 rule %s still exist from sdk", ruleId)
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
		return errors.New("delete CC policy fail, CC policy still exist from sdk")
	}
}
