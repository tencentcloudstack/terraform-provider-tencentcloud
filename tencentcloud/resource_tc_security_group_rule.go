/*
Provides a resource to create security group rule.

Example Usage

Source is CIDR ip

```hcl
resource "tencentcloud_security_group" "sglab_1" {
  name        = "mysg_1"
  description = "favourite sg_1"
  project_id  = 0
}

resource "tencentcloud_security_group_rule" "sglab_1" {
  security_group_id = tencentcloud_security_group.sglab_1.id
  type              = "ingress"
  cidr_ip           = "10.0.0.0/16"
  ip_protocol       = "TCP"
  port_range        = "80"
  policy            = "ACCEPT"
  description       = "favourite sg rule_1"
}
```

Source is a security group id

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
  security_group_id = tencentcloud_security_group.sglab_2.id
  type              = "ingress"
  ip_protocol       = "TCP"
  port_range        = "80"
  policy            = "ACCEPT"
  source_sgid       = tencentcloud_security_group.sglab_3.id
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

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Description:  "Type of the security group rule. Valid values: `ingress` and `egress`.",
			},
			"cidr_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ConflictsWith: []string{
					"source_sgid",
				},
				//ValidateFunc: func(v interface{}, k string) (ws []string, errs []error) {
				//	if _, err := validateIp(v, k); len(err) == 0 {
				//		return
				//	}
				//
				//	if _, err := validateCIDRNetworkAddress(v, k); len(err) != 0 {
				//		errs = append(errs, fmt.Errorf("%s %v is not valid IP address or valid CIDR IP address",
				//			k, v))
				//	}
				//	return
				//},
				Description: "An IP address network or segment, and conflict with `source_sgid` and `address_template`.",
			},
			"ip_protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ConflictsWith: []string{
					"protocol_template",
				},
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validateAllowedStringValueIgnoreCase([]string{"TCP", "UDP", "ICMP"}),
				Description:  "Type of IP protocol. Valid values: `TCP`, `UDP` and `ICMP`. Default to all types protocol, and conflicts with `protocol_template`.",
			},
			"port_range": {
				Type:     schema.TypeString,
				Optional: true,
				ConflictsWith: []string{
					"protocol_template",
				},
				ForceNew:    true,
				Computed:    true,
				Description: "Range of the port. The available value can be one, multiple or one segment. E.g. `80`, `80,90` and `80-90`. Default to all ports, and confilicts with `protocol_template`.",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					match, _ := regexp.MatchString("^(\\d{1,5},)*\\d{1,5}$|^\\d{1,5}-\\d{1,5}$", value)
					if !match {
						errors = append(errors, fmt.Errorf("%s example: `53`, `80,443` and `80-90`, Not configured to represent all ports", k))
					}
					return
				},
			},
			"policy": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValueIgnoreCase([]string{"ACCEPT", "DROP"}),
				Description:  "Rule policy of security group. Valid values: `ACCEPT` and `DROP`.",
			},
			"source_sgid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ConflictsWith: []string{
					"cidr_ip", "address_template",
				},
				Computed:    true,
				Description: "ID of the nested security group, and conflicts with `cidr_ip` and `address_template`.",
			},
			"address_template": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template_id": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Address template ID, conflicts with `group_id`.",
						},
						"group_id": {
							Type:        schema.TypeString,
							ForceNew:    true,
							Optional:    true,
							Description: "Address template group ID, conflicts with `template_id`.",
						},
					},
				},
				Optional: true,
				ForceNew: true,
				ConflictsWith: []string{
					"cidr_ip", "source_sgid",
				},
				Computed:    true,
				Description: "ID of the address template, and confilicts with `source_sgid` and `cidr_ip`.",
			},
			"protocol_template": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template_id": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Address template ID, conflicts with `group_id`.",
						},
						"group_id": {
							Type:        schema.TypeString,
							ForceNew:    true,
							Optional:    true,
							Description: "Address template group ID, conflicts with `template_id`.",
						},
					},
				},
				Optional: true,
				ForceNew: true,
				ConflictsWith: []string{
					"ip_protocol", "port_range",
				},
				Computed:    true,
				Description: "ID of the address template, and conflict with `ip_protocol`, `port_range`.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Description of the security group rule.",
			},
		},
	}
}

func resourceTencentCloudSecurityGroupRuleCreate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_security_group_rule.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	sgId := d.Get("security_group_id").(string)

	policyType := d.Get("type").(string)

	var (
		cidrIp                  *string
		sourceSgId              *string
		protocol                *string
		portRange               *string
		desc                    *string
		addressTemplateId       *string
		addressTemplateGroupId  *string
		protocolTemplateId      *string
		protocolTemplateGroupId *string
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

	if raw, ok := d.GetOk("address_template"); ok {
		addressRaw := raw.([]interface{})
		address := addressRaw[0].(map[string]interface{})
		if v, ok := address["template_id"]; ok {
			addressTemplateId = common.StringPtr(v.(string))
		}
		if v, ok := address["group_id"]; ok {
			addressTemplateGroupId = common.StringPtr(v.(string))
		}
		//complex conflict check
		if addressTemplateGroupId != nil && *addressTemplateGroupId != "" && addressTemplateId != nil && *addressTemplateId != "" {
			return fmt.Errorf("`address_template.group_id` conflicts with `address_template.template_id`")
		}
	}

	if raw, ok := d.GetOk("protocol_template"); ok {
		protocolRaw := raw.([]interface{})
		protocol := protocolRaw[0].(map[string]interface{})
		if v, ok := protocol["template_id"]; ok {
			protocolTemplateId = common.StringPtr(v.(string))
		}
		if v, ok := protocol["group_id"]; ok {
			protocolTemplateGroupId = common.StringPtr(v.(string))
		}
		//complex conflict check
		if protocolTemplateGroupId != nil && *protocolTemplateGroupId != "" && protocolTemplateId != nil && *protocolTemplateId != "" {
			return fmt.Errorf("`protocol_template.group_id` conflicts with `protocol_template.template_id`")
		}
	}

	desc = common.StringPtr(d.Get("description").(string))

	if cidrIp == nil && sourceSgId == nil && addressTemplateId == nil && addressTemplateGroupId == nil {
		return errors.New("need cidr_ip or source_sgid or address_template")
	}

	/*
		if protocol == nil && portRange == nil && protocolTemplateGroupId == nil && protocolTemplateId == nil{
			return errors.New("need protocol and port_range  or protocol_template")
		}

	*/
	action := d.Get("policy").(string)

	if protocol != nil {
		if strings.ToUpper(*protocol) == "ICMP" && portRange != nil {
			return errors.New("when ip_protocol is ICMP, can't set port_range")
		}
	}

	info := securityGroupRuleBasicInfo{
		SgId:                    sgId,
		Action:                  action,
		CidrIp:                  cidrIp,
		Protocol:                protocol,
		PortRange:               portRange,
		SourceSgId:              sourceSgId,
		PolicyType:              policyType,
		Description:             desc,
		AddressTemplateId:       addressTemplateId,
		AddressTemplateGroupId:  addressTemplateGroupId,
		ProtocolTemplateId:      protocolTemplateId,
		ProtocolTemplateGroupId: protocolTemplateGroupId,
	}

	ruleId, err := service.CreateSecurityGroupPolicy(ctx, info)
	if err != nil {
		return err
	}

	d.SetId(ruleId)

	return resourceTencentCloudSecurityGroupRuleRead(d, m)
}

func resourceTencentCloudSecurityGroupRuleRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_security_group_rule.read")()
	defer inconsistentCheck(d, m)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	ruleId := d.Id()
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		sgId, policyType, policy, e := service.DescribeSecurityGroupPolicy(ctx, ruleId)
		if e != nil {
			return retryError(e)
		}

		if policy == nil {
			d.SetId("")
			return nil
		}

		_ = d.Set("security_group_id", sgId)

		_ = d.Set("type", policyType)

		if policy.CidrBlock != nil && *policy.CidrBlock != "" {
			_ = d.Set("cidr_ip", *policy.CidrBlock)
		}

		if policy.SecurityGroupId != nil && *policy.SecurityGroupId != "" {
			_ = d.Set("source_sgid", *policy.SecurityGroupId)
		}

		if policy.AddressTemplate != nil && ((policy.AddressTemplate.AddressGroupId != nil && *policy.AddressTemplate.AddressGroupId != "") || (policy.AddressTemplate.AddressId != nil && *policy.AddressTemplate.AddressId != "")) {
			addressTemplates := make([]map[string]string, 0)
			addressTemplate := map[string]string{
				"group_id":    "",
				"template_id": "",
			}
			if policy.AddressTemplate.AddressGroupId != nil && *policy.AddressTemplate.AddressGroupId != "" {
				addressTemplate["group_id"] = *policy.AddressTemplate.AddressGroupId
			}
			if policy.AddressTemplate.AddressId != nil && *policy.AddressTemplate.AddressId != "" {
				addressTemplate["template_id"] = *policy.AddressTemplate.AddressId
			}
			addressTemplates = append(addressTemplates, addressTemplate)
			_ = d.Set("address_template", addressTemplates)
		}

		if policy.ServiceTemplate != nil && ((policy.ServiceTemplate.ServiceGroupId != nil && *policy.ServiceTemplate.ServiceGroupId != "") || (policy.ServiceTemplate.ServiceId != nil && *policy.ServiceTemplate.ServiceId != "")) {
			protocolTemplates := make([]map[string]string, 0)
			protocolTemplate := map[string]string{
				"group_id":    "",
				"template_id": "",
			}
			if policy.ServiceTemplate.ServiceGroupId != nil && *policy.ServiceTemplate.ServiceGroupId != "" {
				protocolTemplate["group_id"] = *policy.ServiceTemplate.ServiceGroupId
			}
			if policy.ServiceTemplate.ServiceId != nil && *policy.ServiceTemplate.ServiceId != "" {
				protocolTemplate["template_id"] = *policy.ServiceTemplate.ServiceId
			}
			protocolTemplates = append(protocolTemplates, protocolTemplate)
			_ = d.Set("protocol_template", protocolTemplates)
		}

		if policy.Protocol != nil {
			inputProtocol := d.Get("ip_protocol").(string)
			if inputProtocol == "" {
				inputProtocol = "ALL"
			}
			_ = d.Set("ip_protocol", inputProtocol)
		}

		if policy.Port != nil && *policy.Port != "" {
			_ = d.Set("port_range", *policy.Port)
		}

		_ = d.Set("policy", d.Get("policy").(string))

		if policy.PolicyDescription != nil {
			_ = d.Set("description", *policy.PolicyDescription)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func resourceTencentCloudSecurityGroupRuleDelete(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_security_group_rule.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	ruleId := d.Id()
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
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
