/*
Provides a resource to create a teo application_proxy_rule

Example Usage

```hcl
resource "tencentcloud_teo_application_proxy_rule" "application_proxy_rule" {
  zone_id = &lt;nil&gt;
  proxy_id = &lt;nil&gt;
    proto = &lt;nil&gt;
  port = &lt;nil&gt;
  origin_type = &lt;nil&gt;
  origin_value = &lt;nil&gt;
  origin_port = &lt;nil&gt;
  status = &lt;nil&gt;
  forward_client_ip = &lt;nil&gt;
  session_persist = &lt;nil&gt;
}
```

Import

teo application_proxy_rule can be imported using the id, e.g.

```
terraform import tencentcloud_teo_application_proxy_rule.application_proxy_rule application_proxy_rule_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
	"time"
)

func resourceTencentCloudTeoApplicationProxyRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoApplicationProxyRuleCreate,
		Read:   resourceTencentCloudTeoApplicationProxyRuleRead,
		Update: resourceTencentCloudTeoApplicationProxyRuleUpdate,
		Delete: resourceTencentCloudTeoApplicationProxyRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Site ID.",
			},

			"proxy_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Proxy ID.",
			},

			"rule_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Rule ID.",
			},

			"proto": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Protocol. Valid values: `TCP`, `UDP`.",
			},

			"port": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Valid values:- port number: `80` means port 80.- port range: `81-90` means port range 81-90.",
			},

			"origin_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Origin server type.- `custom`: Specified origins.- `origins`: An origin group.",
			},

			"origin_value": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Origin server information.When `OriginType` is custom, this field value indicates multiple origin servers in either of the following formats:- `IP`- Domain name.When `OriginType` is origins, it indicates the origin group ID.",
			},

			"origin_port": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Valid values:- port number: `80` means port 80.- port range: `81-90` means port range 81-90.",
			},

			"status": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Status of this application proxy rule. Valid values to set is `online` and `offline`.- `online`: Enable.- `offline`: Disable.- `progress`: Deploying.- `stopping`: Disabling.- `fail`: Deployment/Disabling failed.",
			},

			"forward_client_ip": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Passes the client IP. Default value is OFF.When Proto is TCP, valid values:- `TOA`: Pass the client IP via TOA.- `PPV1`: Pass the client IP via Proxy Protocol V1.- `PPV2`: Pass the client IP via Proxy Protocol V2.- `OFF`: Do not pass the client IP.When Proto=UDP, valid values:- `PPV2`: Pass the client IP via Proxy Protocol V2.- `OFF`: Do not pass the client IP.",
			},

			"session_persist": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Specifies whether to enable session persistence. Default value is false.",
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
		response = teo.NewCreateApplicationProxyRuleResponse()
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

	if v, ok := d.GetOk("origin_port"); ok {
		request.OriginPort = helper.String(v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		request.Status = helper.String(v.(string))
	}

	if v, ok := d.GetOk("forward_client_ip"); ok {
		request.ForwardClientIp = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("session_persist"); ok {
		request.SessionPersist = helper.Bool(v.(bool))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().CreateApplicationProxyRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create teo applicationProxyRule failed, reason:%+v", logId, err)
		return err
	}

	zoneId = *response.Response.ZoneId
	d.SetId(strings.Join([]string{zoneId, proxyId, ruleId}, FILED_SP))

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"online"}, 60*readRetryTimeout, time.Second, service.TeoApplicationProxyRuleStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

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

	applicationProxyRule, err := service.DescribeTeoApplicationProxyRuleById(ctx, zoneId, proxyId, ruleId)
	if err != nil {
		return err
	}

	if applicationProxyRule == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TeoApplicationProxyRule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if applicationProxyRule.ZoneId != nil {
		_ = d.Set("zone_id", applicationProxyRule.ZoneId)
	}

	if applicationProxyRule.ProxyId != nil {
		_ = d.Set("proxy_id", applicationProxyRule.ProxyId)
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

	if applicationProxyRule.OriginPort != nil {
		_ = d.Set("origin_port", applicationProxyRule.OriginPort)
	}

	if applicationProxyRule.Status != nil {
		_ = d.Set("status", applicationProxyRule.Status)
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

	var (
		modifyApplicationProxyRuleRequest  = teo.NewModifyApplicationProxyRuleRequest()
		modifyApplicationProxyRuleResponse = teo.NewModifyApplicationProxyRuleResponse()
	)

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

	immutableArgs := []string{"zone_id", "proxy_id", "rule_id", "proto", "port", "origin_type", "origin_value", "origin_port", "status", "forward_client_ip", "session_persist"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("proto") {
		if v, ok := d.GetOk("proto"); ok {
			request.Proto = helper.String(v.(string))
		}
	}

	if d.HasChange("port") {
		if v, ok := d.GetOk("port"); ok {
			portSet := v.(*schema.Set).List()
			for i := range portSet {
				port := portSet[i].(string)
				request.Port = append(request.Port, &port)
			}
		}
	}

	if d.HasChange("origin_type") {
		if v, ok := d.GetOk("origin_type"); ok {
			request.OriginType = helper.String(v.(string))
		}
	}

	if d.HasChange("origin_value") {
		if v, ok := d.GetOk("origin_value"); ok {
			originValueSet := v.(*schema.Set).List()
			for i := range originValueSet {
				originValue := originValueSet[i].(string)
				request.OriginValue = append(request.OriginValue, &originValue)
			}
		}
	}

	if d.HasChange("status") {
		if v, ok := d.GetOk("status"); ok {
			request.Status = helper.String(v.(string))
		}
	}

	if d.HasChange("forward_client_ip") {
		if v, ok := d.GetOk("forward_client_ip"); ok {
			request.ForwardClientIp = helper.String(v.(string))
		}
	}

	if d.HasChange("session_persist") {
		if v, ok := d.GetOkExists("session_persist"); ok {
			request.SessionPersist = helper.Bool(v.(bool))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyApplicationProxyRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update teo applicationProxyRule failed, reason:%+v", logId, err)
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
