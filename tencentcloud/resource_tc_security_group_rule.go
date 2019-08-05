/*
Provides a resource to create security group rule.

Example Usage

```hcl
resource "tencentcloud_security_group" "sglab_1" {
  name        = "mysg_1"
  description = "favourite sg_1"
  project_id  = 0
}
resource "tencentcloud_security_group_rule" "sglab_1" {
  security_group_id = "${tencentcloud_security_group.sglab_1.id}"
  type              = "ingress"
  cidr_ip           = "10.0.0.0/16"
  ip_protocol       = "TCP"
  port_range        = "80"
  policy            = "ACCEPT"
  description       = "favourite sg rule_1"
}
```

```hcl
resource "tencentcloud_security_group" "sglab_2" {
  name        = "mysg_2"
  description = "favourite sg_2"
  project_id  = 0
}
resource "tencentcloud_security_group" "sglab_3" {
  name        = "mysg_3"
  description = "favourite sg_3"
  project_id  = 0
}
resource "tencentcloud_security_group_rule" "sglab_2" {
  security_group_id = "${tencentcloud_security_group.sglab_2.id}"
  type              = "ingress"
  ip_protocol       = "TCP"
  port_range        = "80"
  policy            = "ACCEPT"
  source_sgid       = "${tencentcloud_security_group.sglab_3.id}"
  description       = "favourite sg rule_2"
}
```
*/
package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
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
				ValidateFunc: validateAllowedStringValueIgnoreCase([]string{"ingress", "egress"}),
				Description:  "Type of the security group rule, the available value include `ingress` and `egress`.",
			},
			"cidr_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ConflictsWith: []string{
					"source_sgid",
				},
				ValidateFunc: func(v interface{}, k string) (ws []string, errs []error) {
					if _, err := validateIp(v, k); len(err) == 0 {
						return
					}

					if _, err := validateCIDRNetworkAddress(v, k); len(err) != 0 {
						errs = append(errs, fmt.Errorf("%s %v is not valid IP address or valid CIDR IP address",
							k, v))
					}
					return
				},
				Description: "An IP address network or segment, and conflict with `source_sgid`.",
			},
			"ip_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validateAllowedStringValueIgnoreCase([]string{"TCP", "UDP", "ICMP"}),
				Description:  "Type of ip protocol, the available value include `TCP`, `UDP` and `ICMP`. Default to all types protocol.",
			},
			"port_range": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Range of the port. The available value can be one, multiple or one segment. E.g. `80`, `80,90` and `80-90`. Default to all ports.",
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
				ValidateFunc: validateAllowedStringValueIgnoreCase([]string{"ACCEPT", "DROP"}),
				Description:  "Rule policy of security group, the available value include `ACCEPT` and `DROP`.",
			},
			"source_sgid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ConflictsWith: []string{
					"cidr_ip",
				},
				Computed:    true,
				Description: "ID of the nested security group, and conflict with `cidr_ip`.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "",
				Description: "Description of the security group rule.",
			},
		},
	}
}

func resourceTencentCloudSecurityGroupRuleCreate(d *schema.ResourceData, m interface{}) error {
	defer LogElapsed("resource.tencentcloud_security_group_rule.create")()

	logId := GetLogId(nil)
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

	desc = common.StringPtr(d.Get("description").(string))

	if cidrIp == nil && sourceSgId == nil {
		return errors.New("need cidr_ip or source_sgid")
	}

	action := d.Get("policy").(string)

	if protocol != nil {
		if strings.ToUpper(*protocol) == "ICMP" && portRange != nil {
			return errors.New("when ip_protocol is ICMP, can't set port_range")
		}
	}

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
	defer LogElapsed("resource.tencentcloud_security_group_rule.read")()

	logId := GetLogId(nil)
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

	d.Set("security_group_id", sgId)

	d.Set("type", policyType)

	if policy.CidrBlock != nil && *policy.CidrBlock != "" {
		d.Set("cidr_ip", *policy.CidrBlock)
	}

	if policy.SecurityGroupId != nil && *policy.SecurityGroupId != "" {
		d.Set("source_sgid", *policy.SecurityGroupId)
	}

	if policy.Protocol != nil {
		inputProtocol := d.Get("ip_protocol").(string)
		if inputProtocol == "" {
			inputProtocol = "ALL"
		}
		d.Set("ip_protocol", inputProtocol)
	}

	if policy.Port != nil {
		d.Set("port_range", *policy.Port)
	}

	d.Set("policy", d.Get("policy").(string))

	if policy.PolicyDescription != nil {
		d.Set("description", *policy.PolicyDescription)
	}

	return nil
}

func resourceTencentCloudSecurityGroupRuleDelete(d *schema.ResourceData, m interface{}) error {
	defer LogElapsed("resource.tencentcloud_security_group_rule.delete")()

	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	ruleId := d.Id()
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		e := service.DeleteSecurityGroupPolicy(ctx, ruleId)
		if e != nil {
			return resource.RetryableError(fmt.Errorf("security group delete failed: %s", e.Error()))
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s security group rule delete failed: %s\n ", logId, err.Error())
		return err
	}

	return nil
}
