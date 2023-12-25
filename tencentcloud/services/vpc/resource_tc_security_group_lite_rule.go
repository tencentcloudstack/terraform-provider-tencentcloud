package vpc

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSecurityGroupLiteRule() *schema.Resource {
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
				Description: "Ingress rules set. A rule must match the following format: [action]#[source]#[port]#[protocol]. The available value of 'action' is `ACCEPT` and `DROP`. The 'source' can be an IP address network, segment, security group ID and Address Template ID. The 'port' valid format is `80`, `80,443`, `80-90` or `ALL`. The available value of 'protocol' is `TCP`, `UDP`, `ICMP`, `ALL` and `ppm(g?)-xxxxxxxx`. When 'protocol' is `ICMP` or `ALL`, the 'port' must be `ALL`.",
			},
			"egress": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Egress rules set. A rule must match the following format: [action]#[source]#[port]#[protocol]. The available value of 'action' is `ACCEPT` and `DROP`. The 'source' can be an IP address network, segment, security group ID and Address Template ID. The 'port' valid format is `80`, `80,443`, `80-90` or `ALL`. The available value of 'protocol' is `TCP`, `UDP`, `ICMP`, `ALL` and `ppm(g?)-xxxxxxxx`. When 'protocol' is `ICMP` or `ALL`, the 'port' must be `ALL`.",
			},
		},
	}
}

func resourceTencentCloudSecurityGroupLiteRuleCreate(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_security_group_lite_rule.create")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

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
	defer tccommon.LogElapsed("resource.tencentcloud_security_group_lite_rule.read")()
	defer tccommon.InconsistentCheck(d, m)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()

	service := VpcService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

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
	defer tccommon.LogElapsed("resource.tencentcloud_security_group_lite_rule.update")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()

	service := VpcService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

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

		ingress = nil
	}

	if deleteEgress {
		if err := service.DeleteLiteRules(ctx, id, egress, false); err != nil {
			return err
		}

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
	defer tccommon.LogElapsed("resource.tencentcloud_security_group_lite_rule.delete")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()

	service := VpcService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	return service.DetachAllLiteRulesFromSecurityGroup(ctx, id)
}
