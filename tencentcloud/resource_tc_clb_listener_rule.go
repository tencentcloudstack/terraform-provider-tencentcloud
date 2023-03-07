/*
Provides a resource to create a CLB listener rule.

-> **NOTE:** This resource only be applied to the HTTP or HTTPS listeners.

Example Usage

```hcl
resource "tencentcloud_clb_listener_rule" "foo" {
  listener_id                = "lbl-hh141sn9"
  clb_id                     = "lb-k2zjp9lv"
  domain                     = "foo.net"
  url                        = "/bar"
  health_check_switch        = true
  health_check_interval_time = 5
  health_check_health_num    = 3
  health_check_unhealth_num  = 3
  health_check_http_code     = 2
  health_check_http_path     = "Default Path"
  health_check_http_domain   = "Default Domain"
  health_check_http_method   = "GET"
  certificate_ssl_mode       = "MUTUAL"
  certificate_id             = "VjANRdz8"
  certificate_ca_id          = "VfqO4zkB"
  session_expire_time        = 30
  scheduler                  = "WRR"
}
```
Import

CLB listener rule can be imported using the id (version >= 1.47.0), e.g.

```
$ terraform import tencentcloud_clb_listener_rule.foo lb-7a0t6zqb#lbl-hh141sn9#loc-agg236ys
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
	"github.com/pkg/errors"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudClbListenerRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbListenerRuleCreate,
		Read:   resourceTencentCloudClbListenerRuleRead,
		Update: resourceTencentCloudClbListenerRuleUpdate,
		Delete: resourceTencentCloudClbListenerRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of CLB listener.",
			},
			"clb_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of CLB instance.",
			},
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain name of the listener rule.",
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Url of the listener rule.",
			},
			"health_check_switch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates whether health check is enabled.",
			},
			"health_check_interval_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIntegerInRange(5, 300),
				Description:  "Interval time of health check. Valid value ranges: (5~300) sec. and the default is `5` sec. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in `tencentcloud_clb_listener_rule`.",
			},
			"health_check_health_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIntegerInRange(2, 10),
				Description:  "Health threshold of health check, and the default is `3`. If a success result is returned for the health check 3 consecutive times, indicates that the forwarding is normal. The value range is [2-10]. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in `tencentcloud_clb_listener_rule`.",
			},
			"health_check_unhealth_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIntegerInRange(2, 10),
				Description:  "Unhealthy threshold of health check, and the default is `3`. If the unhealthy result is returned 3 consecutive times, indicates that the forwarding is abnormal. The value range is [2-10].  NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in `tencentcloud_clb_listener_rule`.",
			},
			"health_check_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateAllowedStringValue(HEALTH_CHECK_TYPE),
				Description:  "Type of health check. Valid value is `CUSTOM`, `TCP`, `HTTP`.",
			},
			"health_check_time_out": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIntegerInRange(2, 60),
				Description:  "Time out of health check. The value range is [2-60](SEC).",
			},
			"health_check_http_code": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIntegerInRange(1, 31),
				Description:  "HTTP Status Code. The default is 31. Valid value ranges: [1~31]. `1 means the return value '1xx' is health. `2` means the return value '2xx' is health. `4` means the return value '3xx' is health. `8` means the return value '4xx' is health. 16 means the return value '5xx' is health. If you want multiple return codes to indicate health, need to add the corresponding values. NOTES: The 'HTTP' health check of the 'TCP' listener only supports specifying one health check status code. NOTES: Only supports listeners of 'HTTP' and 'HTTPS' protocol.",
			},
			"health_check_http_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Path of health check. NOTES: Only supports listeners of `HTTP` and `HTTPS` protocol.",
			},
			"health_check_http_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Domain name of health check. NOTES: Only supports listeners of `HTTP` and `HTTPS` protocol.",
			},
			"health_check_http_method": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateAllowedStringValue(CLB_HTTP_METHOD),
				Description:  "Methods of health check. NOTES: Only supports listeners of `HTTP` and `HTTPS` protocol. The default is `HEAD`, the available value are `HEAD` and `GET`.",
			},
			"certificate_ssl_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(CERT_SSL_MODE),
				Description:  "Type of certificate. Valid values: `UNIDIRECTIONAL`, `MUTUAL`. NOTES: Only supports listeners of HTTPS protocol.",
			},
			"certificate_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the server certificate. NOTES: Only supports listeners of HTTPS protocol.",
			},
			"certificate_ca_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the client certificate. NOTES: Only supports listeners of HTTPS protocol.",
			},
			"session_expire_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntegerInRange(30, 3600),
				Description:  "Time of session persistence within the CLB listener. NOTES: Available when scheduler is specified as `WRR`, and not available when listener protocol is `TCP_SSL`.  NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in `tencentcloud_clb_listener_rule`.",
			},
			"http2_switch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicate to apply HTTP2.0 protocol or not.",
			},
			"scheduler": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      CLB_LISTENER_SCHEDULER_WRR,
				ValidateFunc: validateAllowedStringValue(CLB_LISTENER_SCHEDULER),
				Description:  "Scheduling method of the CLB listener rules. Valid values: `WRR`, `IP HASH`, `LEAST_CONN`. The default is `WRR`.  NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in `tencentcloud_clb_listener_rule`.",
			},
			"target_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      CLB_TARGET_TYPE_NODE,
				ValidateFunc: validateAllowedStringValue([]string{CLB_TARGET_TYPE_NODE, CLB_TARGET_TYPE_TARGETGROUP}),
				Description:  "Backend target type. Valid values: `NODE`, `TARGETGROUP`. `NODE` means to bind ordinary nodes, `TARGETGROUP` means to bind target group.",
			},
			"forward_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateAllowedStringValue([]string{"HTTP", "HTTPS", "TRPC"}),
				Description:  "Forwarding protocol between the CLB instance and real server. Valid values: `HTTP`, `HTTPS`, `TRPC`.",
			},
			//computed
			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of this CLB listener rule.",
			},
		},
	}
}

func resourceTencentCloudClbListenerRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_listener_rule.create")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	listenerId := d.Get("listener_id").(string)
	checkErr := ListenerIdCheck(listenerId)
	if checkErr != nil {
		return checkErr
	}
	clbId := d.Get("clb_id").(string)
	protocol := ""
	//get listener protocol
	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		instance, e := clbService.DescribeListenerById(ctx, listenerId, clbId)
		if e != nil {
			return retryError(e)
		}
		if instance != nil {
			protocol = *(instance.Protocol)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s get CLB listener failed, reason:%+v", logId, err)
		return err
	}

	if !(protocol == CLB_LISTENER_PROTOCOL_HTTP || protocol == CLB_LISTENER_PROTOCOL_HTTPS) {
		return fmt.Errorf("[CHECK][CLB listener rule][Create] check: The rule can only be created/modified with listeners of protocol HTTP/HTTPS")
	}
	request := clb.NewCreateRuleRequest()
	request.LoadBalancerId = helper.String(clbId)
	request.ListenerId = helper.String(listenerId)

	//rule set
	var rule clb.RuleInput

	domain := d.Get("domain").(string)
	rule.Domain = helper.String(domain)
	url := d.Get("url").(string)
	rule.Url = helper.String(url)
	rule.TargetType = helper.String(d.Get("target_type").(string))
	if v, ok := d.GetOk("forward_type"); ok {
		rule.ForwardType = helper.String(v.(string))
	}
	scheduler := ""
	if v, ok := d.GetOk("scheduler"); ok {
		if !(protocol == CLB_LISTENER_PROTOCOL_HTTP || protocol == CLB_LISTENER_PROTOCOL_HTTPS) {
			return fmt.Errorf("[CHECK][CLB listener rule][Create] check: Scheduler can only be set with listener protocol TCP/UDP/TCP_SSL or rule of listener HTTP/HTTPS")
		}

		scheduler = v.(string)
		rule.Scheduler = helper.String(scheduler)
	}

	if v, ok := d.GetOk("session_expire_time"); ok {
		if !(protocol == CLB_LISTENER_PROTOCOL_HTTP || protocol == CLB_LISTENER_PROTOCOL_HTTPS) {
			return fmt.Errorf("[CHECK][CLB listener rule][Create] check: session_expire_time can only be set with protocol TCP/UDP or rule of listener HTTP/HTTPS")
		}
		if scheduler != CLB_LISTENER_SCHEDULER_WRR && scheduler != "" {
			return fmt.Errorf("[CHECK][CLB listener rule][Create] check: session_expire_time can only be set when scheduler is WRR ")
		}
		vv := int64(v.(int))
		rule.SessionExpireTime = &vv
	}
	healthSetFlag, healthCheck, healthErr := checkHealthCheckPara(ctx, d, protocol, HEALTH_APPLY_TYPE_RULE)
	if healthErr != nil {
		return healthErr
	}
	if healthSetFlag {
		rule.HealthCheck = healthCheck
	}

	certificateSetFlag, certificateInput, certErr := checkCertificateInputPara(ctx, d, meta)
	if certErr != nil {
		return certErr
	}
	if certificateSetFlag {
		if !(protocol == CLB_LISTENER_PROTOCOL_HTTPS) {
			return fmt.Errorf("[CHECK][CLB listener rule][Create] check: certificate para can only be set with rule of linstener with protocol 'HTTPS'")
		}
		rule.Certificate = certificateInput
	}

	request.Rules = []*clb.RuleInput{&rule}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		requestId := ""
		response, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().CreateRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
			requestId = *response.Response.RequestId
			retryErr := waitForTaskFinish(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
			if retryErr != nil {
				return resource.NonRetryableError(errors.WithStack(retryErr))
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create CLB listener rule failed, reason:%+v", logId, err)
		return err
	}

	locationId := ""
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ruleInstance, ruleErr := clbService.DescribeRuleByPara(ctx, clbId, listenerId, domain, url)
		if ruleErr != nil {
			return retryError(errors.WithStack(ruleErr))
		}
		locationId = *ruleInstance.LocationId
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CLB listener rule failed, reason:%+v", logId, err)
		return err
	}

	//this ID style changes since terraform 1.47.0
	d.SetId(clbId + FILED_SP + listenerId + FILED_SP + locationId)

	// set http2
	if v, ok := d.GetOkExists("http2_switch"); ok {
		http2Switch := v.(bool)
		domainRequest := clb.NewModifyDomainAttributesRequest()
		domainRequest.Http2 = &http2Switch
		domainRequest.LoadBalancerId = &clbId
		domainRequest.ListenerId = &listenerId
		domainRequest.Domain = &domain
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			response, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().ModifyDomainAttributes(domainRequest)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
				requestId := *response.Response.RequestId
				retryErr := waitForTaskFinish(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
				if retryErr != nil {
					return resource.NonRetryableError(errors.WithStack(retryErr))
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update CLB listener rule failed, reason:%+v", logId, err)
			return err
		}
	}
	return resourceTencentCloudClbListenerRuleRead(d, meta)
}

func resourceTencentCloudClbListenerRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_listener_rule.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	resourceId := d.Id()
	var locationId = resourceId
	items := strings.Split(resourceId, FILED_SP)
	itemLength := len(items)
	clbId := d.Get("clb_id").(string)
	listenerId := d.Get("listener_id").(string)
	checkErr := ListenerIdCheck(listenerId)
	if checkErr != nil {
		return checkErr
	}
	if itemLength == 1 && clbId == "" {
		return fmt.Errorf("The old style listenerId %s does not support import, please use clbId#listenerId style", resourceId)
	} else if itemLength == 3 && clbId == "" {
		locationId = items[2]
		listenerId = items[1]
		clbId = items[0]
	} else if itemLength == 3 && clbId != "" {
		locationId = items[2]
		listenerId = items[1]
	} else if itemLength != 1 && itemLength != 3 {
		return fmt.Errorf("broken ID %s", resourceId)
	}

	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	//this function is not supported by api, need to be travelled
	filter := map[string]string{"rule_id": locationId, "listener_id": listenerId, "clb_id": clbId}
	var instances []*clb.RuleOutput
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := clbService.DescribeRulesByFilter(ctx, filter)
		if e != nil {
			return retryError(e)
		}
		instances = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CLB listener rule failed, reason:%+v", logId, err)
		return err
	}

	if len(instances) == 0 {
		d.SetId("")
		return nil
	}

	instance := instances[0]
	_ = d.Set("clb_id", clbId)
	_ = d.Set("listener_id", listenerId)
	_ = d.Set("domain", instance.Domain)
	_ = d.Set("rule_id", instance.LocationId)
	_ = d.Set("url", instance.Url)
	_ = d.Set("scheduler", instance.Scheduler)
	_ = d.Set("session_expire_time", instance.SessionExpireTime)
	_ = d.Set("target_type", instance.TargetType)
	_ = d.Set("forward_type", instance.ForwardType)
	_ = d.Set("http2_switch", instance.Http2)

	//health check
	if instance.HealthCheck != nil {
		health_check_switch := false
		if *instance.HealthCheck.HealthSwitch == int64(1) {
			health_check_switch = true
		}
		_ = d.Set("health_check_switch", health_check_switch)
		_ = d.Set("health_check_interval_time", instance.HealthCheck.IntervalTime)
		_ = d.Set("health_check_health_num", instance.HealthCheck.HealthNum)
		_ = d.Set("health_check_unhealth_num", instance.HealthCheck.UnHealthNum)
		_ = d.Set("health_check_http_method", helper.String(strings.ToUpper(*instance.HealthCheck.HttpCheckMethod)))
		_ = d.Set("health_check_http_domain", instance.HealthCheck.HttpCheckDomain)
		_ = d.Set("health_check_http_path", instance.HealthCheck.HttpCheckPath)
		_ = d.Set("health_check_http_code", instance.HealthCheck.HttpCode)
		_ = d.Set("health_check_type", instance.HealthCheck.CheckType)
		_ = d.Set("health_check_time_out", instance.HealthCheck.TimeOut)
	}

	if instance.Certificate != nil {
		_ = d.Set("certificate_ssl_mode", instance.Certificate.SSLMode)
		_ = d.Set("certificate_id", instance.Certificate.CertId)
		if instance.Certificate.CertCaId != nil {
			_ = d.Set("certificate_ca_id", instance.Certificate.CertCaId)
		}
	}

	return nil
}

func resourceTencentCloudClbListenerRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_listener_rule.update")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	resourceId := d.Id()
	items := strings.Split(resourceId, FILED_SP)
	itemLength := len(items)
	locationId := items[itemLength-1]
	listenerId := d.Get("listener_id").(string)
	checkErr := ListenerIdCheck(listenerId)
	if checkErr != nil {
		return checkErr
	}
	clbId := d.Get("clb_id").(string)
	protocol := ""
	//get listener protocol
	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		instance, e := clbService.DescribeListenerById(ctx, listenerId, clbId)
		if e != nil {
			return retryError(e)
		}
		protocol = *(instance.Protocol)
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s get CLB listener failed, reason:%s\n ", logId, err.Error())
		return err
	}

	changed := false
	url := ""
	scheduler := ""
	sessionExpireTime := 0

	request := clb.NewModifyRuleRequest()
	request.ListenerId = helper.String(listenerId)
	request.LoadBalancerId = helper.String(clbId)
	request.LocationId = helper.String(locationId)
	if d.HasChange("url") {
		changed = true
		url = d.Get("url").(string)
		request.Url = helper.String(url)
	}

	if d.HasChange("forward_type") {
		changed = true
		request.ForwardType = helper.String(d.Get("forward_type").(string))
	}

	if d.HasChange("scheduler") {
		changed = true
		scheduler = d.Get("scheduler").(string)
		if !(protocol == CLB_LISTENER_PROTOCOL_HTTP || protocol == CLB_LISTENER_PROTOCOL_HTTPS) {
			return fmt.Errorf("[CHECK][CLB listener rule %s][Update] check: Scheduler can only be set with listener protocol TCP/UDP/TCP_SSL or rule of listener HTTP/HTTPS", locationId)
		}
		request.Scheduler = helper.String(scheduler)
	}

	if d.HasChange("session_expire_time") {
		changed = true
		sessionExpireTime = d.Get("session_expire_time").(int)
		if !(protocol == CLB_LISTENER_PROTOCOL_HTTP || protocol == CLB_LISTENER_PROTOCOL_HTTPS) {
			return fmt.Errorf("[CHECK][CLB listener rule %s][Update] check: session_expire_time can only be set with protocol TCP/UDP or rule of listener HTTP/HTTPS", locationId)
		}
		if scheduler != CLB_LISTENER_SCHEDULER_WRR && scheduler != "" {
			return fmt.Errorf("[CHECK][CLB listener rule %s][Update] check: session_expire_time can only be set when scheduler is WRR", locationId)
		}
		sessionExpireTime64 := int64(sessionExpireTime)
		request.SessionExpireTime = &sessionExpireTime64
	}

	healthSetFlag, healthCheck, healthErr := checkHealthCheckPara(ctx, d, protocol, HEALTH_APPLY_TYPE_RULE)
	if healthErr != nil {
		return healthErr
	}

	if healthSetFlag {
		changed = true
		request.HealthCheck = healthCheck
	}

	if changed {
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			response, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().ModifyRule(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
				requestId := *response.Response.RequestId
				retryErr := waitForTaskFinish(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
				if retryErr != nil {
					return resource.NonRetryableError(errors.WithStack(retryErr))
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update CLB listener rule failed, reason:%+v", logId, err)
			return err
		}
	}

	//modify domain and ssl
	domainChanged := false
	domainRequest := clb.NewModifyDomainAttributesRequest()
	if d.HasChange("domain") {
		old, new := d.GetChange("domain")
		domainChanged = true
		domainRequest.Domain = helper.String(old.(string))
		domainRequest.NewDomain = helper.String(new.(string))
	} else {
		domainRequest.Domain = helper.String(d.Get("domain").(string))
	}

	if d.HasChange("certificate_id") || d.HasChange("certificate_ca_id ") || d.HasChange("certificate_ssl_mode") {
		domainChanged = true
		certificateSetFlag, certificateInput, certErr := checkCertificateInputPara(ctx, d, meta)
		if certErr != nil {
			return certErr
		}
		if certificateSetFlag {
			if !(protocol == CLB_LISTENER_PROTOCOL_HTTPS) {
				return fmt.Errorf("[CHECK][CLB listener rule][Create] check: certificate para can only be set with rule of linstener with protocol 'HTTPS'")
			}
			domainRequest.Certificate = certificateInput
		}
	}

	if d.HasChange("http2_switch") {
		if v, ok := d.GetOkExists("http2_switch"); ok {
			if !(protocol == CLB_LISTENER_PROTOCOL_HTTPS) {
				return fmt.Errorf("[CHECK][CLB listener rule][Create] check: certificate para can only be set with rule of linstener with protocol 'HTTPS'")
			}
			domainChanged = true
			domainRequest.Http2 = helper.Bool(v.(bool))
		}
	}

	if domainChanged {
		domainRequest.ListenerId = &listenerId
		domainRequest.LoadBalancerId = &clbId
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			response, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().ModifyDomainAttributes(domainRequest)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
				requestId := *response.Response.RequestId
				retryErr := waitForTaskFinish(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
				if retryErr != nil {
					return resource.NonRetryableError(errors.WithStack(retryErr))
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update CLB listener rule failed, reason:%+v", logId, err)
			return err
		}
	}

	return nil
}

func resourceTencentCloudClbListenerRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_listener_rule.delete")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	resourceId := d.Id()
	items := strings.Split(resourceId, FILED_SP)
	itemLength := len(items)
	locationId := items[itemLength-1]
	listenerId := d.Get("listener_id").(string)
	checkErr := ListenerIdCheck(listenerId)
	if checkErr != nil {
		return checkErr
	}
	clbId := d.Get("clb_id").(string)

	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := clbService.DeleteRuleById(ctx, clbId, listenerId, locationId)
		if e != nil {
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete CLB listener rule failed, reason:%+v", logId, err)
		return err
	}
	return nil
}
