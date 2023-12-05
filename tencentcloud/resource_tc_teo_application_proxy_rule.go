package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
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
				Description: "Site ID.",
			},

			"proxy_id": {
				Type:        schema.TypeString,
				Required:    true,
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
				Description: "Valid values: `80` means port 80; `81-90` means port range 81-90.",
			},

			"origin_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Origin server type. Valid values: `custom`: Specified origins; `origins`: An origin group.",
			},

			"origin_port": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Origin port, supported formats: single port: 80; Port segment: 81-90, 81 to 90 ports.",
			},

			"origin_value": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "Origin site information: When `OriginType` is `custom`, it indicates one or more origin sites, such as `['8.8.8.8', '9.9.9.9']` or `OriginValue=['test.com']`; When `OriginType` is `origins`, there is required to be one and only one element, representing the origin site group ID, such as `['origin-537f5b41-162a-11ed-abaa-525400c5da15']`.",
			},

			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Status, the values are: `online`: enabled; `offline`: deactivated; `progress`: being deployed; `stopping`: being deactivated; `fail`: deployment failure/deactivation failure.",
			},

			"forward_client_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Passes the client IP. Default value is `OFF`. When Proto is TCP, valid values: `TOA`: Pass the client IP via TOA; `PPV1`: Pass the client IP via Proxy Protocol V1; `PPV2`: Pass the client IP via Proxy Protocol V2; `OFF`: Do not pass the client IP. When Proto=UDP, valid values: `PPV2`: Pass the client IP via Proxy Protocol V2; `OFF`: Do not pass the client IP.",
			},

			"session_persist": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
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

	if v, ok := d.GetOk("port"); ok {
		portSet := v.(*schema.Set).List()
		for i := range portSet {
			port := portSet[i].(string)
			request.Port = append(request.Port, &port)
		}
	}

	if v, ok := d.GetOk("proto"); ok {
		request.Proto = helper.String(v.(string))
	}

	if v, ok := d.GetOk("origin_type"); ok {
		request.OriginType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("origin_port"); ok {
		request.OriginPort = helper.String(v.(string))
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
			return retryError(e, InternalError)
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

	if response.Response.RuleId != nil {
		ruleId = *response.Response.RuleId
	}

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	err = resource.Retry(60*readRetryTimeout, func() *resource.RetryError {
		instance, errRet := service.DescribeTeoApplicationProxyRule(ctx, zoneId, proxyId, ruleId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if *instance.Status == "online" {
			return nil
		}
		if *instance.Status == "fail" {
			return resource.NonRetryableError(fmt.Errorf("applicationProxyRule status is %v, operate failed.", *instance.Status))
		}
		return resource.RetryableError(fmt.Errorf("applicationProxyRule status is %v, retry...", *instance.Status))
	})
	if err != nil {
		return err
	}

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

	proxyRule, err := service.DescribeTeoApplicationProxyRule(ctx, zoneId, proxyId, ruleId)

	if err != nil {
		return err
	}

	if proxyRule == nil {
		d.SetId("")
		return fmt.Errorf("resource `applicationProxyRule` %s does not exist", ruleId)
	}

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("proxy_id", proxyId)
	_ = d.Set("rule_id", ruleId)

	if proxyRule.Proto != nil {
		_ = d.Set("proto", proxyRule.Proto)
	}

	if proxyRule.Port != nil {
		_ = d.Set("port", proxyRule.Port)
	}

	if proxyRule.OriginType != nil {
		_ = d.Set("origin_type", proxyRule.OriginType)
	}

	if proxyRule.OriginPort != nil {
		_ = d.Set("origin_port", proxyRule.OriginPort)
	}

	if proxyRule.OriginValue != nil {
		_ = d.Set("origin_value", proxyRule.OriginValue)
	}

	if proxyRule.Status != nil {
		_ = d.Set("status", proxyRule.Status)
	}

	if proxyRule.ForwardClientIp != nil {
		_ = d.Set("forward_client_ip", proxyRule.ForwardClientIp)
	}

	if proxyRule.SessionPersist != nil {
		_ = d.Set("session_persist", proxyRule.SessionPersist)
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

	if v, ok := d.GetOk("origin_port"); ok {
		request.OriginPort = helper.String(v.(string))
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
		log.Printf("[CRITAL]%s create teo applicationProxyRule failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("status") {
		if v, ok := d.GetOk("status"); ok {
			statusRequest := teo.NewModifyApplicationProxyRuleStatusRequest()

			statusRequest.ZoneId = &zoneId
			statusRequest.ProxyId = &proxyId
			statusRequest.RuleId = &ruleId
			statusRequest.Status = helper.String(v.(string))

			statusErr := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				statusResult, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyApplicationProxyRuleStatus(statusRequest)
				if e != nil {
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
						logId, request.GetAction(), request.ToJsonString(), statusResult.ToJsonString())
				}
				return nil
			})

			if statusErr != nil {
				log.Printf("[CRITAL]%s create teo applicationProxy failed, reason:%+v", logId, statusErr)
				return statusErr
			}
			_ = d.Set("status", v.(string))
		}
	}

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	err = resource.Retry(60*readRetryTimeout, func() *resource.RetryError {
		instance, errRet := service.DescribeTeoApplicationProxyRule(ctx, zoneId, proxyId, ruleId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if *instance.Status == "online" {
			return nil
		}
		if *instance.Status == "fail" {
			return resource.NonRetryableError(fmt.Errorf("applicationProxyRule status is %v, operate failed.", *instance.Status))
		}
		return resource.RetryableError(fmt.Errorf("applicationProxyRule status is %v, retry...", *instance.Status))
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

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		instance, errRet := service.DescribeTeoApplicationProxyRule(ctx, zoneId, proxyId, ruleId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if *instance.Status == "offline" {
			return nil
		}
		if *instance.Status == "stopping" {
			return resource.RetryableError(fmt.Errorf("applicationProxyRule stopping"))
		}

		statusRequest := teo.NewModifyApplicationProxyRuleStatusRequest()
		statusRequest.ZoneId = &zoneId
		statusRequest.ProxyId = &proxyId
		statusRequest.RuleId = &ruleId
		statusRequest.Status = helper.String("offline")
		_, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyApplicationProxyRuleStatus(statusRequest)
		if e != nil {
			return resource.NonRetryableError(fmt.Errorf("setting applicationProxyRule `status` to offline failed, reason: %v", e))
		}
		return resource.RetryableError(fmt.Errorf("setting applicationProxyRule `status` to offline"))
	})
	if err != nil {
		return err
	}

	if err = service.DeleteTeoApplicationProxyRuleById(ctx, zoneId, proxyId, ruleId); err != nil {
		return err
	}

	return nil
}
