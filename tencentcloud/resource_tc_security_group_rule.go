/*
Provide a resource to create security group rule.

Example Usage

```hcl
resource "tencentcloud_security_group_rule" "sglab" {
  security_group_id = "sg-fh48e762"
  type              = "ingress"
  cidr_ip           = "10.0.0.0/16"
  ip_protocol       = "TCP"
  port_range        = "80"
  policy            = "ACCEPT"
  source_sgid       = "sg-fh48e762"
  description       = "favourite sg rule"
}
```
*/
package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

func resourceTencentCloudSecurityGroupRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSecurityGroupRuleCreate,
		Read:   resourceTencentCloudSecurityGroupRuleRead,
		Delete: resourceTencentCloudSecurityGroupRuleDelete,

		Schema: map[string]*schema.Schema{
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the security group to be queried.",
			},

			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue([]string{"ingress", "egress"}),
				Description:  "Type of the security group rule, the available value include 'ingress' and 'egress'.",
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
				Description: "An IP address network or segment, and conflict with 'source_sgid'.",
			},

			"ip_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue([]string{"TCP", "UDP", "ICMP"}),
				Description:  "Type of ip protocol, the available value include 'TCP', 'UDP' and 'ICMP'. Default to all types protocol.",
			},

			"port_range": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Range of the port. The available value can be one, multiple or one segment. E.g. '80', '80,90' and '80-90'. Default to all ports.",
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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue([]string{"ACCEPT", "DROP"}),
				Description:  "Rule policy of security group, the available value include 'ACCEPT' and 'DROP'.",
			},

			"source_sgid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ConflictsWith: []string{
					"cidr_ip",
				},
				Description: "ID of the nested security group, and conflict with 'cidr_ip'.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Description of the security group rule.",
			},
		},
	}
}

func resourceTencentCloudSecurityGroupRuleCreate(d *schema.ResourceData, m interface{}) error {
	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_security_group_rule.create")()

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

	if cidrIp == nil && sourceSgId == nil {
		return errors.New("need cidr_ip or source_sgid")
	}

	action := d.Get("policy").(string)

	info := securityGroupRuleBasicInfo{
		SgId:        sgId,
		Action:      action,
		CidrIp:      cidrIp,
		Protocol:    protocol,
		PortRange:   portRange,
		SourceSgId:  sourceSgId,
		PolicyType:  policyType,
		Description: desc,
	}

	ruleId, err := service.CreateSecurityGroupPolicy(ctx, info)
	if err != nil {
		return err
	}

	d.SetId(ruleId)

	return resourceTencentCloudSecurityGroupRuleRead(d, m)
}

func resourceTencentCloudSecurityGroupRuleRead(d *schema.ResourceData, m interface{}) error {
	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_security_group_rule.read")()

	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	ruleId := d.Id()

	sgId, policyType, policy, err := service.DescribeSecurityGroupPolicy(ctx, ruleId)
	if err != nil {
		return err
	}

	if policy == nil {
		d.SetId("")
		return nil
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
		_ = d.Set("ip_protocol", strings.ToUpper(*policy.Protocol))
	}

	if policy.Port != nil {
		_ = d.Set("port_range", *policy.Port)
	}

	_ = d.Set("policy", *policy.Action)

	if policy.PolicyDescription != nil {
		_ = d.Set("description", *policy.PolicyDescription)
	}

	return nil
}

func resourceTencentCloudSecurityGroupRuleDelete(d *schema.ResourceData, m interface{}) error {
	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_security_group_rule.delete")()

	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	ruleId := d.Id()

	return service.DeleteSecurityGroupPolicy(ctx, ruleId)
}
