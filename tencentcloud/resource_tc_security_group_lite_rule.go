/*
Provide a resource to create security group some lite rules quickly.

-> **NOTE:** It can't be used with tencentcloud_security_group_rule, and don't create multiple tencentcloud_security_group_rule resources, otherwise it may cause problems.

Example Usage

```hcl
resource "tencentcloud_security_group" "foo" {
  name = "ci-temp-test-sg"
}

resource "tencentcloud_security_group_lite_rule" "foo" {
  security_group_id = tencentcloud_security_group.foo.id

  ingress = [
    "ACCEPT#192.168.1.0/24#80#TCP",
    "DROP#8.8.8.8#80,90#UDP",
    "ACCEPT#0.0.0.0/0#80-90#TCP",
  ]

  egress = [
    "ACCEPT#192.168.0.0/16#ALL#TCP",
    "ACCEPT#10.0.0.0/8#ALL#ICMP",
    "DROP#0.0.0.0/0#ALL#ALL",
  ]
}
```

Import

Security group lite rule can be imported using the id, e.g.

```
  $ terraform import tencentcloud_security_group_lite_rule.foo sg-ey3wmiz1
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSecurityGroupLiteRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSecurityGroupLiteRuleCreate,
		Read:   resourceTencentCloudSecurityGroupLiteRuleRead,
		Update: resourceTencentCloudSecurityGroupLiteRuleUpdate,
		Delete: resourceTencentCloudSecurityGroupLiteRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the security group.",
			},
			"ingress": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Ingress rules set. A rule must match the following format: [action]#[cidr_ip]#[port]#[protocol]. The available value of 'action' is `ACCEPT` and `DROP`. The 'cidr_ip' must be an IP address network or segment. The 'port' valid format is `80`, `80,443`, `80-90` or `ALL`. The available value of 'protocol' is `TCP`, `UDP`, `ICMP` and `ALL`. When 'protocol' is `ICMP` or `ALL`, the 'port' must be `ALL`.",
			},
			"egress": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Egress rules set. A rule must match the following format: [action]#[cidr_ip]#[port]#[protocol]. The available value of 'action' is `ACCEPT` and `DROP`. The 'cidr_ip' must be an IP address network or segment. The 'port' valid format is `80`, `80,443`, `80-90` or `ALL`. The available value of 'protocol' is `TCP`, `UDP`, `ICMP` and `ALL`. When 'protocol' is `ICMP` or `ALL`, the 'port' must be `ALL`.",
			},
		},
	}
}

func resourceTencentCloudSecurityGroupLiteRuleCreate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_security_group_lite_rule.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	sgId := d.Get("security_group_id").(string)

	var (
		ingress []VpcSecurityGroupLiteRule
		egress  []VpcSecurityGroupLiteRule
	)

	if raw, ok := d.GetOk("ingress"); ok {
		ingressStrs := helper.InterfacesStrings(raw.([]interface{}))
		for _, ingressStr := range ingressStrs {
			liteRule, err := parseRule(ingressStr)
			if err != nil {
				return err
			}
			ingress = append(ingress, liteRule)
		}
	}

	if raw, ok := d.GetOk("egress"); ok {
		egressStrs := helper.InterfacesStrings(raw.([]interface{}))
		for _, egressStr := range egressStrs {
			liteRule, err := parseRule(egressStr)
			if err != nil {
				return err
			}
			egress = append(egress, liteRule)
		}
	}

	if err := service.AttachLiteRulesToSecurityGroup(ctx, sgId, ingress, egress); err != nil {
		return err
	}

	d.SetId(sgId)

	return resourceTencentCloudSecurityGroupLiteRuleRead(d, m)
}

func resourceTencentCloudSecurityGroupLiteRuleRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_security_group_lite_rule.read")()
	defer inconsistentCheck(d, m)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()

	service := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	respIngress, respEgress, exist, err := service.DescribeSecurityGroupPolices(ctx, id)
	if err != nil {
		return err
	}

	if !exist {
		d.SetId("")
		return nil
	}

	ingress := make([]string, 0, len(respIngress))
	for _, in := range respIngress {
		ingress = append(ingress, in.String())
	}

	egress := make([]string, 0, len(respEgress))
	for _, eg := range respEgress {
		egress = append(egress, eg.String())
	}

	_ = d.Set("security_group_id", id)
	_ = d.Set("ingress", ingress)
	_ = d.Set("egress", egress)

	return nil
}

func resourceTencentCloudSecurityGroupLiteRuleUpdate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_security_group_lite_rule.update")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()

	service := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	var (
		ingress       []VpcSecurityGroupLiteRule
		egress        []VpcSecurityGroupLiteRule
		deleteIngress bool
		deleteEgress  bool
	)

	if d.HasChange("ingress") {
		if raw, ok := d.GetOk("ingress"); ok {
			ingressStrs := helper.InterfacesStrings(raw.([]interface{}))
			for _, ingressStr := range ingressStrs {
				liteRule, err := parseRule(ingressStr)
				if err != nil {
					return err
				}
				ingress = append(ingress, liteRule)
			}
		} else {
			old, _ := d.GetChange("ingress")
			ingressStrs := helper.InterfacesStrings(old.([]interface{}))
			for _, ingressStr := range ingressStrs {
				liteRule, err := parseRule(ingressStr)
				if err != nil {
					return err
				}
				ingress = append(ingress, liteRule)
			}

			deleteIngress = true
		}
	}

	if d.HasChange("egress") {
		if raw, ok := d.GetOk("egress"); ok {
			egressStrs := helper.InterfacesStrings(raw.([]interface{}))
			for _, egressStr := range egressStrs {
				liteRule, err := parseRule(egressStr)
				if err != nil {
					return err
				}
				egress = append(egress, liteRule)
			}
		} else {
			old, _ := d.GetChange("egress")
			egressStrs := helper.InterfacesStrings(old.([]interface{}))
			for _, egressStr := range egressStrs {
				liteRule, err := parseRule(egressStr)
				if err != nil {
					return err
				}
				egress = append(egress, liteRule)
			}

			deleteEgress = true
		}
	}

	d.Partial(true)

	if deleteIngress && deleteEgress {
		if err := service.DetachAllLiteRulesFromSecurityGroup(ctx, id); err != nil {
			return err
		}

		d.Partial(false)

		return resourceTencentCloudSecurityGroupLiteRuleRead(d, m)
	}

	if deleteIngress {
		if err := service.DeleteLiteRules(ctx, id, ingress, true); err != nil {
			return err
		}

		d.SetPartial("ingress")

		ingress = nil
	}

	if deleteEgress {
		if err := service.DeleteLiteRules(ctx, id, egress, false); err != nil {
			return err
		}

		d.SetPartial("egress")

		egress = nil
	}

	// if both len == 0, means both rules are deleted
	if len(ingress) > 0 || len(egress) > 0 {
		if err := service.modifyLiteRulesInSecurityGroup(ctx, id, ingress, egress); err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceTencentCloudSecurityGroupLiteRuleRead(d, m)
}

func resourceTencentCloudSecurityGroupLiteRuleDelete(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_security_group_lite_rule.delete")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()

	service := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	return service.DetachAllLiteRulesFromSecurityGroup(ctx, id)
}
