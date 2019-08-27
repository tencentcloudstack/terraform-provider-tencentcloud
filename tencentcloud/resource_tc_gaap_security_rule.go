package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTencentCloudGaapSecurityRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGaapSecurityRuleCreate,
		Read:   resourceTencentCloudGaapSecurityRuleRead,
		Update: resourceTencentCloudGaapSecurityRuleUpdate,
		Delete: resourceTencentCloudGaapSecurityRuleDelete,
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cidr_ip": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateCidrIp,
			},
			"action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{"ACCEPT", "DROP"}),
				ForceNew:     true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				ValidateFunc: validateStringLengthInRange(0, 30),
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ALL",
				ValidateFunc: validateAllowedStringValue([]string{"ALL", "TCP", "UDP"}),
				ForceNew:     true,
			},
			"port": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ALL",
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value == "ALL" {
						return
					}
					if !regexp.MustCompile("^(\\d{1,5},)*\\d{1,5}$|^\\d{1,5}-\\d{1,5}$").MatchString(value) {
						errors = append(errors, fmt.Errorf("%s example: 53、80,443、80-90, Not configured to represent all ports", k))
					}
					return
				},
			},
		},
	}
}

func resourceTencentCloudGaapSecurityRuleCreate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_security_rule.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	policyId := d.Get("policy_id").(string)
	cidrIp := d.Get("cidr_ip").(string)
	action := d.Get("action").(string)
	port := d.Get("port").(string)
	protocol := d.Get("protocol").(string)
	name := d.Get("name").(string)

	if protocol == "ALL" && port != "ALL" {
		return errors.New("when protocol is ALL, port should be ALL, too")
	}

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	id, err := service.CreateSecurityRule(ctx, policyId, name, cidrIp, port, action, protocol)
	if err != nil {
		return err
	}

	d.SetId(id)

	return resourceTencentCloudGaapSecurityRuleRead(d, m)
}

func resourceTencentCloudGaapSecurityRuleRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_security_rule.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()
	policyId := d.Get("policy_id").(string)

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	rule, err := service.DescribeSecurityRule(ctx, policyId, id)
	if err != nil {
		return err
	}

	if rule == nil {
		d.SetId("")
		return nil
	}

	if rule.SourceCidr == nil {
		return errors.New("security rule cidr IP is nil")
	}
	d.Set("cidr_ip", rule.SourceCidr)

	if rule.Action == nil {
		return errors.New("security rule action is nil")
	}
	d.Set("action", rule.Action)

	if rule.AliasName == nil {
		return errors.New("security rule name is nil")
	}
	d.Set("name", rule.AliasName)

	if rule.Protocol == nil {
		return errors.New("security rule protocol is nil")
	}
	d.Set("protocol", rule.Protocol)

	if rule.DestPortRange == nil {
		return errors.New("security rule port is nil")
	}
	d.Set("port", rule.DestPortRange)

	return nil
}

func resourceTencentCloudGaapSecurityRuleUpdate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_security_rule.update")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()
	policyId := d.Get("policy_id").(string)
	name := d.Get("name").(string)

	if name == "" {
		return errors.New("new name can't be empty")
	}

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	if err := service.ModifySecurityRuleName(ctx, policyId, id, name); err != nil {
		return err
	}

	return resourceTencentCloudGaapSecurityRuleRead(d, m)
}

func resourceTencentCloudGaapSecurityRuleDelete(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_security_rule.delete")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()
	policyId := d.Get("policy_id").(string)

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	return service.DeleteSecurityRule(ctx, policyId, id)
}
