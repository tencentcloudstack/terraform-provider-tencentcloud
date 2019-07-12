package tencentcloud

import (
	"context"
	"fmt"
	"net"
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

func resourceTencentCloudSecurityGroupRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSecurityGroupRuleCreate,
		Read:   resourceTencentCloudSecurityGroupRuleRead,
		Update: resourceTencentCloudSecurityGroupRuleUpdate,
		Delete: resourceTencentCloudSecurityGroupRuleDelete,

		Schema: map[string]*schema.Schema{
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value != "ingress" && value != "egress" {
						errors = append(errors, fmt.Errorf("%s of rule, ingress (inbound) or egress (outbound) value:%v", k, value))
					}
					return
				},
			},

			"cidr_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ConflictsWith: []string{
					"source_sgid",
				},
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					if _, _, err := net.ParseCIDR(v.(string)); err != nil {
						errors = append(errors, err)
					}

					return
				},
			},

			"ip_protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					switch value {
					case "TCP", "UDP", "ICMP":

					default:
						errors = append(errors, fmt.Errorf("%s support 'UDP', 'TCP', 'ICMP' and not configured means all protocols. But got %s", k, v))
					}
					return
				},
			},

			"port_range": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "example: 53、80,443、80-90",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					match, _ := regexp.MatchString("^(\\d{1,5},)*\\d{1,5}$|^\\d{1,5}\\-\\d{1,5}$", value)
					if !match {
						errors = append(errors, fmt.Errorf("%s example: 53、80,443、80-90, Not configured to represent all ports", k))
					}
					return
				},
			},

			"policy": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value != "ACCEPT" && value != "DROP" {
						errors = append(errors, fmt.Errorf("policy of rule, 'ACCEPT' or 'DROP'"))
					}
					return
				},
			},

			"source_sgid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ConflictsWith: []string{
					"cidr_ip",
				},
			},

			"policy_index": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Description:  "security group rule index, start from 0",
				ValidateFunc: validateIntegerMin(0),
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "security group rule description",
			},
		},
	}
}

func resourceTencentCloudSecurityGroupRuleCreate(d *schema.ResourceData, m interface{}) error {
	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_security_group.read")()

	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	sgId := d.Get("security_group_id").(string)

	policyType := d.Get("type").(string)

	var (
		cidrIp     *string
		sourceSgId *string
		protocol   *string
		portRange  *string
		desc       *string
		index      *int64
	)

	if raw, ok := d.GetOk("cidr_ip"); ok {
		cidrIp = common.StringPtr(raw.(string))
	}

	if raw, ok := d.GetOk("source_sgid"); ok {
		sourceSgId = common.StringPtr(raw.(string))
	}

	if raw, ok := d.GetOk("ip_protocol"); ok {
		protocol = common.StringPtr(raw.(string))
	}

	if raw, ok := d.GetOk("port_range"); ok {
		portRange = common.StringPtr(raw.(string))
	}

	if raw, ok := d.GetOk("description"); ok {
		desc = common.StringPtr(raw.(string))
	}

	if raw, ok := d.GetOk("policy_index"); ok {
		index = common.Int64Ptr(int64(raw.(int)))
	}

	action := d.Get("policy").(string)

	ruleId, err := service.CreateSecurityGroupPolicy(
		ctx,
		sgId,
		policyType,
		action,
		cidrIp,
		sourceSgId,
		protocol,
		portRange,
		desc,
		index,
	)
	if err != nil {
		return err
	}

	d.SetId(ruleId)

	return resourceTencentCloudSecurityGroupRuleRead(d, m)
}

func resourceTencentCloudSecurityGroupRuleRead(d *schema.ResourceData, m interface{}) error {
	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_security_group.read")()

	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	ruleId := d.Id()

	sgId, policyType, policy, err := service.DescribeSecurityGroupPolicy(ctx, ruleId)
	if err != nil {
		return err
	}

	_ = d.Set("security_group_id", sgId)

	_ = d.Set("type", policyType)

	if policy.CidrBlock != nil {
		_ = d.Set("cidr_ip", *policy.CidrBlock)
	}

	if policy.SecurityGroupId != nil {
		_ = d.Set("source_sgid", *policy.SecurityGroupId)
	}

	if policy.Protocol != nil {
		_ = d.Set("ip_protocol", *policy.Protocol)
	}

	if policy.Port != nil {
		_ = d.Set("port_range", *policy.Port)
	}

	_ = d.Set("policy", *policy.Action)

	_ = d.Set("policy_index", int(*policy.PolicyIndex))

	if policy.PolicyDescription != nil {
		_ = d.Set("description", *policy.PolicyDescription)
	}

	return nil
}

func resourceTencentCloudSecurityGroupRuleUpdate(d *schema.ResourceData, m interface{}) error {
	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_security_group.read")()

	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	ruleId := d.Id()
	index := d.Get("policy_index").(int)

	if d.HasChange("description") {
		var desc *string
		if descRaw, ok := d.GetOk("description"); ok {
			desc = common.StringPtr(descRaw.(string))
		}

		if err := service.ModifySecurityGroupPolicy(ctx, ruleId, index, desc); err != nil {
			return err
		}

		return resourceTencentCloudSecurityGroupRuleRead(d, m)
	}

	return nil
}

func resourceTencentCloudSecurityGroupRuleDelete(d *schema.ResourceData, m interface{}) error {
	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_security_group.read")()

	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	ruleId := d.Id()

	index := d.Get("policy_index").(int)

	return service.DeleteSecurityGroupPolicy(ctx, ruleId, index)
}
