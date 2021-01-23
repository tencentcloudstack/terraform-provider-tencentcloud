/*
Provides a resource to create a security policy rule.

Example Usage

```hcl
resource "tencentcloud_gaap_proxy" "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

resource "tencentcloud_gaap_security_policy" "foo" {
  proxy_id = tencentcloud_gaap_proxy.foo.id
  action   = "ACCEPT"
}

resource "tencentcloud_gaap_security_rule" "foo" {
  policy_id = tencentcloud_gaap_security_policy.foo.id
  cidr_ip   = "1.1.1.1"
  action    = "ACCEPT"
  protocol  = "TCP"
}
```

Import

GAAP security rule can be imported using the id, e.g.

```
  $ terraform import tencentcloud_gaap_security_rule.foo sr-xxxxxxxx
```
*/
package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceTencentCloudGaapSecurityRule() *schema.Resource {
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
				ForceNew:     true,
				ValidateFunc: validateCidrIp,
				Description:  "A network address block of the request source.",
			},
			"action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{"ACCEPT", "DROP"}),
				ForceNew:     true,
				Description:  "Policy of the rule. Valid value: `ACCEPT` and `DROP`.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				ValidateFunc: validateStringLengthInRange(0, 30),
				Description:  "Name of the security policy rule. Maximum length is 30.",
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ALL",
				ValidateFunc: validateAllowedStringValue([]string{"ALL", "TCP", "UDP"}),
				ForceNew:     true,
				Description:  "Protocol of the security policy rule. Default value is `ALL`. Valid value: `TCP`, `UDP` and `ALL`.",
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
	defer logElapsed("resource.tencentcloud_gaap_security_rule.create")()
	gaapActionMu.Lock()
	defer gaapActionMu.Unlock()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

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
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	rule, err := service.DescribeSecurityRule(ctx, id)
	if err != nil {
		return err
	}

	if rule == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("policy_id", rule.PolicyId)
	_ = d.Set("cidr_ip", rule.SourceCidr)
	_ = d.Set("action", rule.Action)
	_ = d.Set("name", rule.AliasName)
	_ = d.Set("protocol", rule.Protocol)
	_ = d.Set("port", rule.DestPortRange)

	return nil
}

func resourceTencentCloudGaapSecurityRuleUpdate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_security_rule.update")()
	gaapActionMu.Lock()
	defer gaapActionMu.Unlock()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

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
	gaapActionMu.Lock()
	defer gaapActionMu.Unlock()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()
	policyId := d.Get("policy_id").(string)

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	return service.DeleteSecurityRule(ctx, policyId, id)
}
