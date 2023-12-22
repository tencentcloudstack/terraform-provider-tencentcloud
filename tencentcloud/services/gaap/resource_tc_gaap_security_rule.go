package gaap

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudGaapSecurityRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGaapSecurityRuleCreate,
		Read:   resourceTencentCloudGaapSecurityRuleRead,
		Update: resourceTencentCloudGaapSecurityRuleUpdate,
		Delete: resourceTencentCloudGaapSecurityRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the security policy.",
			},
			"cidr_ip": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateCidrIp,
				Description:  "A network address block of the request source.",
			},
			"action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"ACCEPT", "DROP"}),
				Description:  "Policy of the rule. Valid value: `ACCEPT` and `DROP`.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				ValidateFunc: tccommon.ValidateStringLengthInRange(0, 30),
				Description:  "Name of the security policy rule. Maximum length is 30.",
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ALL",
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"ALL", "TCP", "UDP"}),
				Description:  "Protocol of the security policy rule. Default value is `ALL`. Valid value: `TCP`, `UDP` and `ALL`.",
			},
			"port": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ALL",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value == "ALL" {
						return
					}
					if !regexp.MustCompile(`^(\d{1,5},)*\d{1,5}$|^\d{1,5}-\d{1,5}$`).MatchString(value) {
						errors = append(errors, fmt.Errorf("%s example: `53`, `80,443` and `80-90`, Not configured to represent all ports", k))
					}
					return
				},
				Description: "Target port. Default value is `ALL`. Valid examples: `80`, `80,443` and `3306-20000`.",
			},
		},
	}
}

func resourceTencentCloudGaapSecurityRuleCreate(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_security_rule.create")()
	gaapActionMu.Lock()
	defer gaapActionMu.Unlock()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	policyId := d.Get("policy_id").(string)
	cidrIp := d.Get("cidr_ip").(string)
	action := d.Get("action").(string)
	port := d.Get("port").(string)
	protocol := d.Get("protocol").(string)
	name := d.Get("name").(string)

	if protocol == "ALL" && port != "ALL" {
		return errors.New("when protocol is ALL, port should be ALL, too")
	}

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	id, err := service.CreateSecurityRule(ctx, policyId, name, cidrIp, port, action, protocol)
	if err != nil {
		return err
	}

	d.SetId(id)

	return resourceTencentCloudGaapSecurityRuleRead(d, m)
}

func resourceTencentCloudGaapSecurityRuleRead(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_security_rule.read")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	rule, err := service.DescribeSecurityRule(ctx, id)
	if err != nil {
		return err
	}

	if rule == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("policy_id", rule.PolicyId)

	cidrIp := *rule.SourceCidr
	// fix when cidr is "x.x.x.x/32", because return will remove /32
	if v, ok := d.GetOk("cidr_ip"); ok {
		getCidrIp := v.(string)
		splits := strings.Split(getCidrIp, "/")
		if len(splits) > 1 {
			if splits[1] == "32" && cidrIp == splits[0] {
				cidrIp = fmt.Sprintf("%s/32", cidrIp)
			}
		}
	}

	_ = d.Set("cidr_ip", cidrIp)
	_ = d.Set("action", rule.Action)
	_ = d.Set("name", rule.AliasName)
	_ = d.Set("protocol", rule.Protocol)
	_ = d.Set("port", rule.DestPortRange)

	return nil
}

func resourceTencentCloudGaapSecurityRuleUpdate(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_security_rule.update")()
	gaapActionMu.Lock()
	defer gaapActionMu.Unlock()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()
	policyId := d.Get("policy_id").(string)
	cidrIp := d.Get("cidr_ip").(string)
	action := d.Get("action").(string)
	port := d.Get("port").(string)
	protocol := d.Get("protocol").(string)
	name := d.Get("name").(string)

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	if err := service.ModifySecurityRule(ctx, policyId, id, cidrIp, action, port, protocol, name); err != nil {
		return err
	}

	return resourceTencentCloudGaapSecurityRuleRead(d, m)
}

func resourceTencentCloudGaapSecurityRuleDelete(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_security_rule.delete")()
	gaapActionMu.Lock()
	defer gaapActionMu.Unlock()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()
	policyId := d.Get("policy_id").(string)

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	return service.DeleteSecurityRule(ctx, policyId, id)
}
