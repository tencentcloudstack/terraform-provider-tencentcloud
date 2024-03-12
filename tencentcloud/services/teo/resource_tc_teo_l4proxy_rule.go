package teo

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTeoL4proxy_rule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoL4proxy_ruleCreate,
		Read:   resourceTencentCloudTeoL4proxy_ruleRead,
		Update: resourceTencentCloudTeoL4proxy_ruleUpdate,
		Delete: resourceTencentCloudTeoL4proxy_ruleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Zone ID.",
			},

			"proxy_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Layer 4 proxy instance ID.",
			},

			"rule_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Forwarding rule ID.\nNote: Do not fill in this parameter when L4ProxyRule is used as an input parameter in CreateL4ProxyRules; it must be filled in when L4ProxyRule is used as an input parameter in ModifyL4ProxyRules.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Forwarding protocol. Valid values:\n<li>TCP: TCP protocol;</li>\n<li>UDP: UDP protocol.</li>\nNote: This parameter must be filled in when L4ProxyRule is used as an input parameter in CreateL4ProxyRules; it is optional when L4ProxyRule is used as an input parameter in ModifyL4ProxyRules. If not specified, it will retain its existing value.",
			},
			"port_range": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "Forwarding port, which can be set as follows:\n<li>A single port, such as 80;</li>\n<li>A range of ports, such as 81-85, representing ports 81, 82, 83, 84, 85.</li>\nNote: This parameter must be filled in when L4ProxyRule is used as an input parameter in CreateL4ProxyRules; it is optional when L4ProxyRule is used as an input parameter in ModifyL4ProxyRules. If not specified, it will retain its existing value.",
			},
			"origin_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Origin server type. Valid values:\n<li>IP_DOMAIN: IP/Domain name origin server;</li>\n<li>ORIGIN_GROUP: Origin server group;</li>\n<li>LB: Cloud Load Balancer, currently only open to the allowlist.</li>\nNote: This parameter must be filled in when L4ProxyRule is used as an input parameter in CreateL4ProxyRules; it is optional when L4ProxyRule is used as an input parameter in ModifyL4ProxyRules. If not specified, it will retain its existing value.",
			},
			"origin_value": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "Origin server address.\n<li>When OriginType is set to IP_DOMAIN, enter the IP address or domain name, such as 8.8.8.8 or test.com;</li>\n<li>When OriginType is set to ORIGIN_GROUP, enter the origin server group ID, such as og-537y24vf5b41;</li>\n<li>When OriginType is set to LB, enter the Cloud Load Balancer instance ID, such as lb-2qwk30xf7s9g.</li>\nNote: This parameter must be filled in when L4ProxyRule is used as an input parameter in CreateL4ProxyRules; it is optional when L4ProxyRule is used as an input parameter in ModifyL4ProxyRules. If not specified, it will retain its existing value.",
			},
			"origin_port_range": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Origin server port, which can be set as follows:<li>A single port, such as 80;</li>\n<li>A range of ports, such as 81-85, representing ports 81, 82, 83, 84, 85. When inputting a range of ports, ensure that the length corresponds with that of the forwarding port range. For example, if the forwarding port range is 80-90, this port range should be 90-100.</li>\nNote: This parameter must be filled in when L4ProxyRule is used as an input parameter in CreateL4ProxyRules; it is optional when L4ProxyRule is used as an input parameter in ModifyL4ProxyRules. If not specified, it will retain its existing value.",
			},
			"client_ip_pass_through_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Transmission of the client IP address. Valid values:<li>TOA: Available only when Protocol=TCP;</li> \n<li>PPV1: Transmission via Proxy Protocol V1. Available only when Protocol=TCP;</li>\n<li>PPV2: Transmission via Proxy Protocol V2;</li> \n<li>SPP: Transmission via Simple Proxy Protocol. Available only when Protocol=UDP;</li> \n<li>OFF: No transmission.</li>\nNote: This parameter is optional when L4ProxyRule is used as an input parameter in CreateL4ProxyRules, and if not specified, the default value OFF will be used; it is optional when L4ProxyRule is used as an input parameter in ModifyL4ProxyRules. If not specified, it will retain its existing value.",
			},
			"session_persist": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies whether to enable session persistence. Valid values:\n<li>on: Enable;</li>\n<li>off: Disable.</li>\nNote: This parameter is optional when L4ProxyRule is used as an input parameter in CreateL4ProxyRules, and if not specified, the default value off will be used; it is optional when L4ProxyRule is used as an input parameter in ModifyL4ProxyRules. If not specified, it will retain its existing value.",
			},
			"session_persist_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Session persistence period, with a range of 30-3600, measured in seconds.\nNote: This parameter is optional when L4ProxyRule is used as an input parameter in CreateL4ProxyRules. It is valid only when SessionPersist is set to on and defaults to 3600 if not specified. It is optional when L4ProxyRule is used as an input parameter in ModifyL4ProxyRules. If not specified, it will retain its existing value.",
			},
			"rule_tag": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Rule tag. Accepts 1-50 arbitrary characters.\nNote: This parameter is optional when L4ProxyRule is used as an input parameter in CreateL4ProxyRules; it is optional when L4ProxyRule is used as an input parameter in ModifyL4ProxyRules. If not specified, it will retain its existing value.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Rule status. Valid values:<li>online: Enabled;</li>\n<li>offline: Disabled;</li>\n<li>progress: Deploying;</li>\n<li>stopping: Disabling;</li>\n<li>fail: Failed to deploy or disable.</li>\nNote: Do not set this parameter when L4ProxyRule is used as an input parameter in CreateL4ProxyRules and ModifyL4ProxyRules.",
			},
		},
	}
}

func resourceTencentCloudTeoL4proxy_ruleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_l4proxy_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = teo.NewCreateL4ProxyRulesRequest()
		response = teo.NewCreateL4ProxyRulesResponse()
		zoneId   string
		proxyId  string
		ruleId   string
	)
	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
		request.ZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("proxy_id"); ok {
		proxyId = v.(string)
		request.ProxyId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("l4_proxy_rules"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			l4ProxyRule := teo.L4ProxyRule{}
			if v, ok := dMap["rule_id"]; ok {
				l4ProxyRule.RuleId = helper.String(v.(string))
			}
			if v, ok := dMap["protocol"]; ok {
				l4ProxyRule.Protocol = helper.String(v.(string))
			}
			if v, ok := dMap["port_range"]; ok {
				portRangeSet := v.(*schema.Set).List()
				for i := range portRangeSet {
					if portRangeSet[i] != nil {
						portRange := portRangeSet[i].(string)
						l4ProxyRule.PortRange = append(l4ProxyRule.PortRange, &portRange)
					}
				}
			}
			if v, ok := dMap["origin_type"]; ok {
				l4ProxyRule.OriginType = helper.String(v.(string))
			}
			if v, ok := dMap["origin_value"]; ok {
				originValueSet := v.(*schema.Set).List()
				for i := range originValueSet {
					if originValueSet[i] != nil {
						originValue := originValueSet[i].(string)
						l4ProxyRule.OriginValue = append(l4ProxyRule.OriginValue, &originValue)
					}
				}
			}
			if v, ok := dMap["origin_port_range"]; ok {
				l4ProxyRule.OriginPortRange = helper.String(v.(string))
			}
			if v, ok := dMap["client_ip_pass_through_mode"]; ok {
				l4ProxyRule.ClientIPPassThroughMode = helper.String(v.(string))
			}
			if v, ok := dMap["session_persist"]; ok {
				l4ProxyRule.SessionPersist = helper.String(v.(string))
			}
			if v, ok := dMap["session_persist_time"]; ok {
				l4ProxyRule.SessionPersistTime = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["rule_tag"]; ok {
				l4ProxyRule.RuleTag = helper.String(v.(string))
			}
			if v, ok := dMap["status"]; ok {
				l4ProxyRule.Status = helper.String(v.(string))
			}
			request.L4ProxyRules = append(request.L4ProxyRules, &l4ProxyRule)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().CreateL4ProxyRules(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create teo l4proxy failed, reason:%+v", logId, err)
		return err
	}

	ruleId = *response.Response.L4ProxyRuleIds[0]
	d.SetId(strings.Join([]string{zoneId, proxyId, ruleId}, tccommon.FILED_SP))

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	conf := tccommon.BuildStateChangeConf([]string{}, []string{"online"}, 10*tccommon.ReadRetryTimeout, time.Second, service.TeoL4proxyRuleStateRefreshFunc(zoneId, proxyId, ruleId, []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudTeoL4proxy_ruleRead(d, meta)
}

func resourceTencentCloudTeoL4proxy_ruleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_l4proxy_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	proxyId := idSplit[1]
	ruleId := idSplit[2]

	l4proxy_rule, err := service.DescribeTeoL4proxyRuleById(ctx, zoneId, proxyId, ruleId)
	if err != nil {
		return err
	}

	if l4proxy_rule == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TeoL4proxy_rule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("proxy_id", proxyId)
	if l4proxy_rule.RuleId != nil {
		_ = d.Set("rule_id", l4proxy_rule.RuleId)
	}

	if l4proxy_rule.L4ProxyRules != nil {
		l4ProxyRulesList := []interface{}{}
		for _, l4ProxyRules := range l4proxy_rule.L4ProxyRules {
			l4ProxyRulesMap := map[string]interface{}{}

			if l4proxy_rule.L4ProxyRules.RuleId != nil {
				l4ProxyRulesMap["rule_id"] = l4proxy_rule.L4ProxyRules.RuleId
			}

			if l4proxy_rule.L4ProxyRules.Protocol != nil {
				l4ProxyRulesMap["protocol"] = l4proxy_rule.L4ProxyRules.Protocol
			}

			if l4proxy_rule.L4ProxyRules.PortRange != nil {
				l4ProxyRulesMap["port_range"] = l4proxy_rule.L4ProxyRules.PortRange
			}

			if l4proxy_rule.L4ProxyRules.OriginType != nil {
				l4ProxyRulesMap["origin_type"] = l4proxy_rule.L4ProxyRules.OriginType
			}

			if l4proxy_rule.L4ProxyRules.OriginValue != nil {
				l4ProxyRulesMap["origin_value"] = l4proxy_rule.L4ProxyRules.OriginValue
			}

			if l4proxy_rule.L4ProxyRules.OriginPortRange != nil {
				l4ProxyRulesMap["origin_port_range"] = l4proxy_rule.L4ProxyRules.OriginPortRange
			}

			if l4proxy_rule.L4ProxyRules.ClientIPPassThroughMode != nil {
				l4ProxyRulesMap["client_ip_pass_through_mode"] = l4proxy_rule.L4ProxyRules.ClientIPPassThroughMode
			}

			if l4proxy_rule.L4ProxyRules.SessionPersist != nil {
				l4ProxyRulesMap["session_persist"] = l4proxy_rule.L4ProxyRules.SessionPersist
			}

			if l4proxy_rule.L4ProxyRules.SessionPersistTime != nil {
				l4ProxyRulesMap["session_persist_time"] = l4proxy_rule.L4ProxyRules.SessionPersistTime
			}

			if l4proxy_rule.L4ProxyRules.RuleTag != nil {
				l4ProxyRulesMap["rule_tag"] = l4proxy_rule.L4ProxyRules.RuleTag
			}

			if l4proxy_rule.L4ProxyRules.Status != nil {
				l4ProxyRulesMap["status"] = l4proxy_rule.L4ProxyRules.Status
			}

			l4ProxyRulesList = append(l4ProxyRulesList, l4ProxyRulesMap)
		}

		_ = d.Set("l4_proxy_rules", l4ProxyRulesList)

	}

	return nil
}

func resourceTencentCloudTeoL4proxy_ruleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_l4proxy_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = teo.NewModifyL4ProxyRulesRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	proxyId := idSplit[1]
	ruleId := idSplit[2]

	request.ZoneId = &zoneId
	request.ProxyId = &proxyId

	immutableArgs := []string{"zone_id", "proxy_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("l4_proxy_rules") {
		if v, ok := d.GetOk("l4_proxy_rules"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				l4ProxyRule := teo.L4ProxyRule{}

				l4ProxyRule.RuleId = &ruleId

				if v, ok := dMap["protocol"]; ok {
					l4ProxyRule.Protocol = helper.String(v.(string))
				}
				if v, ok := dMap["port_range"]; ok {
					portRangeSet := v.(*schema.Set).List()
					for i := range portRangeSet {
						if portRangeSet[i] != nil {
							portRange := portRangeSet[i].(string)
							l4ProxyRule.PortRange = append(l4ProxyRule.PortRange, &portRange)
						}
					}
				}
				if v, ok := dMap["origin_type"]; ok {
					l4ProxyRule.OriginType = helper.String(v.(string))
				}
				if v, ok := dMap["origin_value"]; ok {
					originValueSet := v.(*schema.Set).List()
					for i := range originValueSet {
						if originValueSet[i] != nil {
							originValue := originValueSet[i].(string)
							l4ProxyRule.OriginValue = append(l4ProxyRule.OriginValue, &originValue)
						}
					}
				}
				if v, ok := dMap["origin_port_range"]; ok {
					l4ProxyRule.OriginPortRange = helper.String(v.(string))
				}
				if v, ok := dMap["client_ip_pass_through_mode"]; ok {
					l4ProxyRule.ClientIPPassThroughMode = helper.String(v.(string))
				}
				if v, ok := dMap["session_persist"]; ok {
					l4ProxyRule.SessionPersist = helper.String(v.(string))
				}
				if v, ok := dMap["session_persist_time"]; ok {
					l4ProxyRule.SessionPersistTime = helper.IntUint64(v.(int))
				}
				if v, ok := dMap["rule_tag"]; ok {
					l4ProxyRule.RuleTag = helper.String(v.(string))
				}
				if v, ok := dMap["status"]; ok {
					l4ProxyRule.Status = helper.String(v.(string))
				}
				request.L4ProxyRules = append(request.L4ProxyRules, &l4ProxyRule)
			}
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().ModifyL4ProxyRules(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update teo l4proxy_rule failed, reason:%+v", logId, err)
		return err
	}

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	conf := tccommon.BuildStateChangeConf([]string{}, []string{"online"}, 10*tccommon.ReadRetryTimeout, time.Second, service.TeoL4proxyRuleStateRefreshFunc(zoneId, proxyId, ruleId, []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudTeoL4proxy_ruleRead(d, meta)
}

func resourceTencentCloudTeoL4proxy_ruleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_l4proxy_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	proxyId := idSplit[1]
	ruleId := idSplit[2]

	if err := service.DeleteTeoL4proxyRuleById(ctx, zoneId, proxyId, ruleId); err != nil {
		return err
	}

	return nil
}
