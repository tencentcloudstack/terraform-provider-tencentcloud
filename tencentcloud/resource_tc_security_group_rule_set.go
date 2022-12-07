/*
Provides a resource to create security group rule.
This resource is similar with tencentcloud_security_group_lite_rule, rules can be ordered and configure descriptions

~> **NOTE:** This resource must exclusive in one security group, do not declare additional rule resources of this security group elsewhere.

Example Usage

```hcl
resource "tencentcloud_security_group" "sglab_1" {
  name        = "mysg_1"
  description = "favourite sg_1"
}

resource "tencentcloud_security_group_rule_set" "sglab_1" {
  security_group_id = tencentcloud_security_group.sglab_1.id
  ingress {
    cidr_block  = "10.0.0.0/16" # Accept IP or CIDR
    protocol    = "TCP" # Default is ALL
    port        = "80" # Accept port e.g. 80 or PortRange e.g. 8080-8089
    action      = "ACCEPT"
    description = "favourite sg rule_1"
  }
  ingress {
    protocol           = "TCP"
    port               = "80"
    action             = "ACCEPT"
    source_security_id = tencentcloud_security_group.sglab_3.id
    description        = "favourite sg rule_2"
  }

  egress {
    action              = "ACCEPT"
    address_template_id = "ipm-xxxxxxxx" # Support address template (group)
    description         = "Allow address template"
  }
  egress {
    action                 = "ACCEPT"
    service_template_group = "ppmg-xxxxxxxx" # Support protocol template (group)
    description            = "Allow protocol template"
  }
  egress {
    cidr_block  = "10.0.0.0/16"
    protocol    = "TCP"
    port        = "80"
    action      = "DROP"
    description = "favourite sg egress rule"
  }
}
```


*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceTencentCloudSecurityGroupRuleSet() *schema.Resource {
	ruleElem := map[string]*schema.Schema{
		"action": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateAllowedStringValueIgnoreCase([]string{"ACCEPT", "DROP"}),
			Description:  "Rule policy of security group. Valid values: `ACCEPT` and `DROP`.",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Description of the security group rule.",
		},
		"cidr_block": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "An IP address network or CIDR segment. NOTE: `cidr_block`, `ipv6_cidr_block`, `source_security_id` and `address_template_*` are exclusive and cannot be set in the same time.",
		},
		"ipv6_cidr_block": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "An IPV6 address network or CIDR segment, and conflict with `source_security_id` and `address_template_*`.",
		},
		"source_security_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "ID of the nested security group, and conflicts with `cidr_block` and `address_template_*`.",
		},
		"address_template_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specify Address template ID like `ipm-xxxxxxxx`, conflict with `source_security_id` and `cidr_block`.",
		},
		"address_template_group": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specify Group ID of Address template like `ipmg-xxxxxxxx`, conflict with `source_security_id` and `cidr_block`.",
		},
		"service_template_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specify Protocol template ID like `ppm-xxxxxxxx`, conflict with `cidr_block` and `port`.",
		},
		"service_template_group": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specify Group ID of Protocol template ID like `ppmg-xxxxxxxx`, conflict with `cidr_block` and `port`.",
		},
		"protocol": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Type of IP protocol. Valid values: `TCP`, `UDP` and `ICMP`. Default to all types protocol, and conflicts with `service_template_*`.",
		},
		"port": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Range of the port. The available value can be one, multiple or one segment. E.g. `80`, `80,90` and `80-90`. Default to all ports, and conflicts with `service_template_*`.",
		},
	}
	return &schema.Resource{
		Create: resourceTencentCloudSecurityGroupRuleSetCreate,
		Read:   resourceTencentCloudSecurityGroupRuleSetRead,
		Update: resourceTencentCloudSecurityGroupRuleSetUpdate,
		Delete: resourceTencentCloudSecurityGroupRuleSetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the security group to be queried.",
			},
			"ingress": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of ingress rule. NOTE: this block is ordered, the first rule has the highest priority.",
				Elem: &schema.Resource{
					Schema: ruleElem,
				},
			},
			"egress": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of egress rule. NOTE: this block is ordered, the first rule has the highest priority.",
				Elem: &schema.Resource{
					Schema: ruleElem,
				},
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Security policies version, auto increment for every update.",
			},
		},
	}
}

func resourceTencentCloudSecurityGroupRuleSetCreate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_security_group_rule_set.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	var err error
	id := d.Get("security_group_id").(string)
	request := vpc.NewModifySecurityGroupPoliciesRequest()
	request.SecurityGroupId = helper.String(id)
	request.SecurityGroupPolicySet = &vpc.SecurityGroupPolicySet{}

	if v, ok := d.GetOk("ingress"); ok {
		rules := v.([]interface{})
		request.SecurityGroupPolicySet.Ingress, err = unmarshalSecurityPolicy(rules)
		if err != nil {
			return err
		}
	}
	if v, ok := d.GetOk("egress"); ok {
		rules := v.([]interface{})
		request.SecurityGroupPolicySet.Egress, err = unmarshalSecurityPolicy(rules)
		if err != nil {
			return err
		}
	}

	err = service.ModifySecurityGroupPolicies(ctx, request)
	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceTencentCloudSecurityGroupRuleSetRead(d, m)
}

func resourceTencentCloudSecurityGroupRuleSetRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_security_group_rule_set.read")()
	defer inconsistentCheck(d, m)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	securityGroupId := d.Id()
	request := vpc.NewDescribeSecurityGroupPoliciesRequest()
	request.SecurityGroupId = &securityGroupId

	result, err := service.DescribeSecurityGroupPolicies(ctx, securityGroupId)
	if err != nil {
		return err
	}

	_ = d.Set("security_group_id", securityGroupId)
	d.SetId(securityGroupId)
	_ = d.Set("version", result.Version)
	if len(result.Ingress) > 0 {
		_ = d.Set("ingress", marshalSecurityPolicy(result.Ingress))
	}
	if len(result.Egress) > 0 {
		_ = d.Set("egress", marshalSecurityPolicy(result.Egress))
	}
	return nil
}

func resourceTencentCloudSecurityGroupRuleSetUpdate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("tencentcloud_security_group_rule_set.update")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	client := m.(*TencentCloudClient).apiV3Conn
	service := VpcService{client}

	version := d.Get("version").(string)
	ver, vErr := strconv.ParseInt(version, 10, 64)
	nextVer := ""

	request := vpc.NewModifySecurityGroupPoliciesRequest()
	request.SecurityGroupId = helper.String(d.Id())
	request.SecurityGroupPolicySet = &vpc.SecurityGroupPolicySet{}
	request.SortPolicys = helper.Bool(true)
	if vErr == nil {
		nextVer = fmt.Sprintf("%d", ver+1)
		request.SecurityGroupPolicySet.Version = &nextVer
	}

	var err error
	if d.HasChange("ingress") {
		rules := d.Get("ingress").([]interface{})
		request.SecurityGroupPolicySet.Ingress, err = unmarshalSecurityPolicy(rules)
		if err != nil {
			return err
		}
	}
	if d.HasChange("egress") {
		rules := d.Get("egress").([]interface{})
		request.SecurityGroupPolicySet.Egress, err = unmarshalSecurityPolicy(rules)
		if err != nil {
			return err
		}
	}
	err = service.ModifySecurityGroupPolicies(ctx, request)
	if err != nil {
		return err
	}

	return resourceTencentCloudSecurityGroupRuleSetRead(d, m)
}

func resourceTencentCloudSecurityGroupRuleSetDelete(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_security_group_rule_set.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	id := d.Id()

	request := vpc.NewModifySecurityGroupPoliciesRequest()
	request.SecurityGroupId = &id
	request.SecurityGroupPolicySet = &vpc.SecurityGroupPolicySet{
		Version: helper.String("0"),
		Ingress: []*vpc.SecurityGroupPolicy{},
		Egress:  []*vpc.SecurityGroupPolicy{},
	}
	//if v, ok := d.GetOk("ingress"); ok {
	//	rules := v.([]interface{})
	//	request.SecurityGroupPolicySet.Ingress, _ = unmarshalSecurityPolicy(rules)
	//}
	//if v, ok := d.GetOk("egress"); ok {
	//	rules := v.([]interface{})
	//	request.SecurityGroupPolicySet.Egress, _ = unmarshalSecurityPolicy(rules)
	//}
	err := service.ModifySecurityGroupPolicies(ctx, request)
	if err != nil {
		log.Printf("[CRITAL]%s security group rule delete failed: %s\n ", logId, err.Error())
		return err
	}

	return nil
}

func unmarshalSecurityPolicy(policies []interface{}) (output []*vpc.SecurityGroupPolicy, err error) {
	for i := range policies {
		policy := policies[i].(map[string]interface{})
		result := &vpc.SecurityGroupPolicy{
			Action: helper.String(policy["action"].(string)),
		}
		// CidrBlock, Ipv6CidrBlock, SecurityGroupId, AddressTemplate are exclusive, and Protocol + Port, ServiceTemplate are also exclusive
		var (
			cidrBlock            = policy["cidr_block"].(string)
			ipv6CidrBlock        = policy["ipv6_cidr_block"].(string)
			sgId                 = policy["source_security_id"].(string)
			addressTemplateId    = policy["address_template_id"].(string)
			addressTemplateGroup = policy["address_template_group"].(string)
			protocol             = policy["protocol"].(string)
			port                 = policy["port"].(string)
			serviceTemplate      = policy["service_template_id"].(string)
			serviceTemplateGroup = policy["service_template_group"].(string)
			desc                 = policy["description"].(string)
		)

		// check if exclusive arguments both set
		checkExcludeValues := func(item map[string]string) (result []string) {
			for k, v := range item {
				if v != "" {
					result = append(result, k)
				}
			}
			return result
		}

		if excludes := checkExcludeValues(map[string]string{
			"cidr_block":             cidrBlock,
			"ipv6_cidr_block":        ipv6CidrBlock,
			"source_security_id":     sgId,
			"address_template_id":    addressTemplateId,
			"address_template_group": addressTemplateGroup,
		}); len(excludes) > 1 {
			err = fmt.Errorf("conflict at rule.%d, cannot set %s in time", i, strings.Join(excludes, ","))
			return
		}

		if excludes := checkExcludeValues(map[string]string{
			"protocol + port":        protocol + port,
			"service_template_id":    serviceTemplate,
			"service_template_group": serviceTemplateGroup,
		}); len(excludes) > 1 {
			err = fmt.Errorf("conflict at rule.%d, cannot set %s in time", i, strings.Join(excludes, ","))
			return
		}

		if cidrBlock != "" {
			result.CidrBlock = &cidrBlock
		}
		if ipv6CidrBlock != "" {
			result.Ipv6CidrBlock = &ipv6CidrBlock
		}
		if sgId != "" {
			result.SecurityGroupId = &sgId
		}
		if addressTemplateId != "" || addressTemplateGroup != "" {
			result.AddressTemplate = &vpc.AddressTemplateSpecification{}
		}
		if addressTemplateId != "" {
			result.AddressTemplate.AddressId = &addressTemplateId
		}
		if addressTemplateGroup != "" {
			result.AddressTemplate.AddressGroupId = &addressTemplateGroup
		}
		if protocol != "" {
			result.Protocol = &protocol
		}
		if port != "" {
			result.Port = &port
		}
		if serviceTemplate != "" || serviceTemplateGroup != "" {
			result.ServiceTemplate = &vpc.ServiceTemplateSpecification{}
		}
		if serviceTemplate != "" {
			result.ServiceTemplate.ServiceId = &serviceTemplate
		}
		if serviceTemplateGroup != "" {
			result.ServiceTemplate.ServiceGroupId = &serviceTemplateGroup
		}
		if desc != "" {
			result.PolicyDescription = &desc
		}
		//result.PolicyIndex = helper.IntInt64(i)

		output = append(output, result)
	}
	return
}

func marshalSecurityPolicy(policies []*vpc.SecurityGroupPolicy) []interface{} {
	result := make([]interface{}, 0, len(policies))
	for i := range policies {
		policy := policies[i]
		dMap := map[string]interface{}{
			"action": policy.Action,
		}
		if policy.CidrBlock != nil {
			dMap["cidr_block"] = policy.CidrBlock
		}
		if policy.Ipv6CidrBlock != nil {
			dMap["ipv6_cidr_block"] = policy.Ipv6CidrBlock
		}
		if policy.Ipv6CidrBlock != nil {
			dMap["source_security_id"] = policy.SecurityGroupId
		}
		if policy.AddressTemplate != nil && policy.AddressTemplate.AddressId != nil {
			dMap["address_template_id"] = policy.AddressTemplate.AddressId
		}
		if policy.AddressTemplate != nil && policy.AddressTemplate.AddressGroupId != nil {
			dMap["address_template_group"] = policy.AddressTemplate.AddressGroupId
		}
		if policy.Protocol != nil /*!checkPolicyPortIgnore(policy.Protocol, policy)*/ {
			dMap["protocol"] = strings.ToUpper(*policy.Protocol)
		}
		if policy.Port != nil /*!checkPolicyPortIgnore(policy.Port, policy)*/ {
			dMap["port"] = policy.Port
		}
		if policy.ServiceTemplate != nil && policy.ServiceTemplate.ServiceId != nil {
			dMap["service_template_id"] = policy.ServiceTemplate.ServiceId
		}
		if policy.ServiceTemplate != nil && policy.ServiceTemplate.ServiceGroupId != nil {
			dMap["service_template_group"] = policy.ServiceTemplate.ServiceGroupId
		}
		if policy.PolicyDescription != nil {
			dMap["description"] = policy.PolicyDescription
		}
		result = append(result, dMap)
	}
	return result
}
