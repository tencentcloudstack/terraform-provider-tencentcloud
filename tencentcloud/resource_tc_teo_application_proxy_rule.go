/*
Provides a resource to create a teo application_proxy_rule

Example Usage

```hcl
resource "tencentcloud_teo_application_proxy_rule" "application_proxy_rule" {
  zone_id  = tencentcloud_teo_zone.zone.id
  proxy_id = tencentcloud_teo_application_proxy.application_proxy_rule.proxy_id

  forward_client_ip = "TOA"
  origin_type       = "custom"
  origin_value      = [
    "1.1.1.1:80",
  ]
  port = [
    "80",
  ]
  proto           = "TCP"
  session_persist = false
}

```
Import

teo application_proxy_rule can be imported using the id, e.g.
```
$ terraform import tencentcloud_teo_application_proxy_rule.application_proxy_rule zoneId#proxyId#ruleId
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220106"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTeoApplicationProxyRule() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTeoApplicationProxyRuleRead,
		Create: resourceTencentCloudTeoApplicationProxyRuleCreate,
		Update: resourceTencentCloudTeoApplicationProxyRuleUpdate,
		Delete: resourceTencentCloudTeoApplicationProxyRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Site ID.",
			},

			"proxy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Proxy ID.",
			},

			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule ID.",
			},

			"proto": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Protocol. Valid values: `TCP`, `UDP`.",
			},

			"port": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "Valid values:- port number: `80` means port 80.- port range: `81-90` means port range 81-90.",
			},

			"origin_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Origin server type.- custom: Specified origins.- origins: An origin group.- load_balancing: A load balancer.",
			},

			"origin_value": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "Origin server information.When OriginType is custom, this field value indicates multiple origin servers in either of the following formats:- IP:Port- Domain name:Port.When OriginType is origins, it indicates the origin group ID.",
			},

			"forward_client_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Passes the client IP.When Proto is TCP, valid values:- TOA: Pass the client IP via TOA.- PPV1: Pass the client IP via Proxy Protocol V1.- PPV2: Pass the client IP via Proxy Protocol V2.- OFF: Do not pass the client IP.When Proto=UDP, valid values:- PPV2: Pass the client IP via Proxy Protocol V2.- OFF: Do not pass the client IP.",
			},

			"session_persist": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Specifies whether to enable session persistence.",
			},
		},
	}
}

func resourceTencentCloudTeoApplicationProxyRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_application_proxy_rule.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = teo.NewCreateApplicationProxyRuleRequest()
		response *teo.CreateApplicationProxyRuleResponse
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

	if v, ok := d.GetOk("proto"); ok {
		request.Proto = helper.String(v.(string))
	}

	if v, ok := d.GetOk("port"); ok {
		portSet := v.(*schema.Set).List()
		for i := range portSet {
			port := portSet[i].(string)
			request.Port = append(request.Port, &port)
		}
	}

	if v, ok := d.GetOk("origin_type"); ok {
		request.OriginType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("origin_value"); ok {
		originValueSet := v.(*schema.Set).List()
		for i := range originValueSet {
			originValue := originValueSet[i].(string)
			request.OriginValue = append(request.OriginValue, &originValue)
		}
	}

	if v, ok := d.GetOk("forward_client_ip"); ok {
		request.ForwardClientIp = helper.String(v.(string))
	}

	if v, ok := d.GetOk("session_persist"); ok {
		request.SessionPersist = helper.Bool(v.(bool))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().CreateApplicationProxyRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create teo applicationProxyRule failed, reason:%+v", logId, err)
		return err
	}

	ruleId = *response.Response.RuleId

	d.SetId(zoneId + FILED_SP + proxyId + FILED_SP + ruleId)
	return resourceTencentCloudTeoApplicationProxyRuleRead(d, meta)
}

func resourceTencentCloudTeoApplicationProxyRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_application_proxy_rule.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	proxyId := idSplit[1]
	ruleId := idSplit[2]

	applicationProxyRule, err := service.DescribeTeoApplicationProxyRule(ctx, zoneId, proxyId, ruleId)

	if err != nil {
		return err
	}

	if applicationProxyRule == nil {
		d.SetId("")
		return fmt.Errorf("resource `applicationProxyRule` %s does not exist", ruleId)
	}

	if applicationProxyRule.RuleId != nil {
		_ = d.Set("rule_id", applicationProxyRule.RuleId)
	}

	if applicationProxyRule.Proto != nil {
		_ = d.Set("proto", applicationProxyRule.Proto)
	}

	if applicationProxyRule.Port != nil {
		_ = d.Set("port", applicationProxyRule.Port)
	}

	if applicationProxyRule.OriginType != nil {
		_ = d.Set("origin_type", applicationProxyRule.OriginType)
	}

	if applicationProxyRule.OriginValue != nil {
		_ = d.Set("origin_value", applicationProxyRule.OriginValue)
	}

	if applicationProxyRule.ForwardClientIp != nil {
		_ = d.Set("forward_client_ip", applicationProxyRule.ForwardClientIp)
	}

	if applicationProxyRule.SessionPersist != nil {
		_ = d.Set("session_persist", applicationProxyRule.SessionPersist)
	}

	return nil
}

func resourceTencentCloudTeoApplicationProxyRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_application_proxy_rule.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := teo.NewModifyApplicationProxyRuleRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	proxyId := idSplit[1]
	ruleId := idSplit[2]

	request.ZoneId = &zoneId
	request.ProxyId = &proxyId
	request.RuleId = &ruleId

	if d.HasChange("zone_id") {
		return fmt.Errorf("`zone_id` do not support change now.")
	}

	if d.HasChange("proxy_id") {
		return fmt.Errorf("`proxy_id` do not support change now.")
	}

	if v, ok := d.GetOk("proto"); ok {
		request.Proto = helper.String(v.(string))
	}

	if v, ok := d.GetOk("port"); ok {
		portSet := v.(*schema.Set).List()
		for i := range portSet {
			port := portSet[i].(string)
			request.Port = append(request.Port, &port)
		}
	}

	if v, ok := d.GetOk("origin_type"); ok {
		request.OriginType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("origin_value"); ok {
		originValueSet := v.(*schema.Set).List()
		for i := range originValueSet {
			originValue := originValueSet[i].(string)
			request.OriginValue = append(request.OriginValue, &originValue)
		}
	}

	if d.HasChange("forward_client_ip") {
		if v, ok := d.GetOk("forward_client_ip"); ok {
			request.ForwardClientIp = helper.String(v.(string))
		}
	}

	if d.HasChange("session_persist") {
		if v, ok := d.GetOk("session_persist"); ok {
			request.SessionPersist = helper.Bool(v.(bool))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyApplicationProxyRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudTeoApplicationProxyRuleRead(d, meta)
}

func resourceTencentCloudTeoApplicationProxyRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_application_proxy_rule.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	proxyId := idSplit[1]
	ruleId := idSplit[2]

	if err := service.DeleteTeoApplicationProxyRuleById(ctx, zoneId, proxyId, ruleId); err != nil {
		return err
	}

	return nil
}
