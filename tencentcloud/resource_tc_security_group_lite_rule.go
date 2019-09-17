package tencentcloud

import (
	"context"
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ingress": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"egress": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceTencentCloudSecurityGroupLiteRuleCreate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_security_group_lite_rule.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	sgId := d.Get("security_group_id").(string)

	var (
		ingress []VpcSecurityGroupLiteRule
		egress  []VpcSecurityGroupLiteRule
	)

	if raw, ok := d.GetOk("ingress"); ok {
		ingressStrs := expandStringList(raw.([]interface{}))
		for _, ingressStr := range ingressStrs {
			action, cidrIp, port, protocol, err := parseRule(ingressStr)
			if err != nil {
				return err
			}
			ingress = append(ingress, VpcSecurityGroupLiteRule{
				action:   action,
				cidrIp:   cidrIp,
				port:     port,
				protocol: protocol,
			})
		}
	}

	if raw, ok := d.GetOk("egress"); ok {
		egressStrs := expandStringList(raw.([]interface{}))
		for _, egressStr := range egressStrs {
			action, cidrIp, port, protocol, err := parseRule(egressStr)
			if err != nil {
				return err
			}
			egress = append(egress, VpcSecurityGroupLiteRule{
				action:   action,
				cidrIp:   cidrIp,
				port:     port,
				protocol: protocol,
			})
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
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

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

	d.Set("security_group_id", id)
	d.Set("ingress", ingress)
	d.Set("egress", egress)

	return nil
}

func resourceTencentCloudSecurityGroupLiteRuleUpdate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_security_group_lite_rule.update")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()

	service := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	var (
		ingress []VpcSecurityGroupLiteRule
		egress  []VpcSecurityGroupLiteRule
	)

	if raw, ok := d.GetOk("ingress"); ok {
		ingressStrs := expandStringList(raw.([]interface{}))
		for _, ingressStr := range ingressStrs {
			action, cidrIp, port, protocol, err := parseRule(ingressStr)
			if err != nil {
				return err
			}
			ingress = append(ingress, VpcSecurityGroupLiteRule{
				action:   action,
				cidrIp:   cidrIp,
				port:     port,
				protocol: protocol,
			})
		}
	}

	if raw, ok := d.GetOk("egress"); ok {
		egressStrs := expandStringList(raw.([]interface{}))
		for _, egressStr := range egressStrs {
			action, cidrIp, port, protocol, err := parseRule(egressStr)
			if err != nil {
				return err
			}
			egress = append(egress, VpcSecurityGroupLiteRule{
				action:   action,
				cidrIp:   cidrIp,
				port:     port,
				protocol: protocol,
			})
		}
	}

	if err := service.AttachLiteRulesToSecurityGroup(ctx, id, ingress, egress); err != nil {
		return err
	}

	return resourceTencentCloudSecurityGroupLiteRuleRead(d, m)
}

func resourceTencentCloudSecurityGroupLiteRuleDelete(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_security_group_lite_rule.delete")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()

	service := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	return service.DetachAllLiteRulesFromSecurityGroup(ctx, id)
}

func parseRule(str string) (action, cidrIp, port, protocol string, err error) {
	split := strings.Split(str, "#")
	if len(split) != 4 {
		err = fmt.Errorf("invalid security group rule %s", str)
		return
	}

	action, cidrIp, port, protocol = split[0], split[1], split[2], split[3]

	switch action {
	default:
		err = fmt.Errorf("invalid action %s, allow action is `ACCEPT` or `DROP`", action)
		return
	case "ACCEPT", "DROP":
	}

	if net.ParseIP(cidrIp) == nil {
		if _, _, err = net.ParseCIDR(cidrIp); err != nil {
			err = fmt.Errorf("invalid cidr_ip %s, allow cidr_ip format is `8.8.8.8` or `10.0.1.0/24`", cidrIp)
			return
		}
	}

	if port != "ALL" && !regexp.MustCompile(`^(\d{1,5},)*\d{1,5}$|^\d{1,5}-\d{1,5}$`).MatchString(port) {
		err = fmt.Errorf("invalid port %s, allow port format is `ALL`, `53`, `80,443` or `80-90`", port)
		return
	}

	switch protocol {
	default:
		err = fmt.Errorf("invalid protocol %s, allow protocol is `ALL`, `TCP`, `UDP` or `ICMP`", protocol)
		return

	case "ALL", "ICMP":
		if port != "ALL" {
			err = fmt.Errorf("when protocol is %s, port must be ALL", protocol)
			return
		}

		// when protocol is ALL or ICMP, port should be "" to avoid sdk error
		port = ""

	case "TCP", "UDP":
	}

	return
}
